package ydocr_test

import (
	"testing"
	"bytes"
	"mytest/ydocr"
	"log"
	"os"
	"io/ioutil"
	"image"
	"fmt"
)

const appKey  = "00eff76a22438df4"
const appSecret  = "5np5oUdN4lz4nNKnCXG3JePiGB51iWO7"
const filename = `C:\Users\GBA\go\src\mytest\ydocr\f1.jpg`

type testpair struct {
	decoded, encoded string
}

var pairs = []testpair{
	// RFC 3548 examples
	{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l+"},
	{"\x14\xfb\x9c\x03\xd9", "FPucA9k="},
	{"\x14\xfb\x9c\x03", "FPucAw=="},

	// RFC 4648 examples
	{"", ""},
	{"f", "Zg=="},
	{"fo", "Zm8="},
	{"foo", "Zm9v"},
	{"foob", "Zm9vYg=="},
	{"fooba", "Zm9vYmE="},
	{"foobar", "Zm9vYmFy"},

	// Wikipedia examples
	{"sure.", "c3VyZS4="},
	{"sure", "c3VyZQ=="},
	{"sur", "c3Vy"},
	{"su", "c3U="},
	{"leasure.", "bGVhc3VyZS4="},
	{"easure.", "ZWFzdXJlLg=="},
	{"asure.", "YXN1cmUu"},
	{"sure.", "c3VyZS4="},
}

type md5Test struct {
	out string
	in  string
}

var golden = []md5Test{
	{"d41d8cd98f00b204e9800998ecf8427e", ""},
	{"0cc175b9c0f1b6a831c399e269772661", "a"},
	{"187ef4436122d1cc2f40dc2b92f0eba0", "ab"},
	{"900150983cd24fb0d6963f7d28e17f72", "abc"},
	{"e2fc714c4727ee9395f324cd2e7f331f", "abcd"},
	{"ab56b4d92b40713acc5af89985d4b786", "abcde"},
	{"e80b5017098950fc58aad83c8c14978e", "abcdef"},
	{"7ac66c0f148de9519b8bd264312c4d64", "abcdefg"},
	{"e8dc4081b13434b45189a720b77b6818", "abcdefgh"},
	{"8aa99b1f439ff71293e95357bac6fd94", "abcdefghi"},
	{"a925576942e94b2ef57a066101b48876", "abcdefghij"},
	{"d747fc1719c7eacb84058196cfe56d57", "Discard medicine more than two years old."},
	{"bff2dcb37ef3a44ba43ab144768ca837", "He who has a shady past knows that nice guys finish last."},
	{"0441015ecb54a7342d017ed1bcfdbea5", "I wouldn't marry him with a ten foot pole."},
	{"9e3cac8e9e9757a60c3ea391130d3689", "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{"a0f04459b031f916a59a35cc482dc039", "The days of the digital watch are numbered.  -Tom Stoppard"},
	{"e7a48e0fe884faf31475d2a04b1362cc", "Nepal premier won't resign."},
	{"637d2fe925c07c113800509964fb0e06", "For every action there is an equal and opposite government program."},
	{"834a8d18d5c6562119cf4c7f5086cb71", "His money is twice tainted: 'taint yours and 'taint mine."},
	{"de3a4d2fd6c73ec2db2abad23b444281", "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{"acf203f997e2cf74ea3aff86985aefaf", "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{"e1c1384cb4d2221dfdd7c795a4222c9a", "size:  a.out:  bad magic"},
	{"c90f3ddecc54f34228c063d7525bf644", "The major problem is with sendmail.  -Mark Horton"},
	{"cdf7ab6c1fd49bd9933c43f3ea5af185", "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{"83bc85234942fc883c063cbd7f0ad5d0", "If the enemy is within range, then so are you."},
	{"277cbe255686b48dd7e8f389394d9299", "It's well we cannot hear the screams/That we create in others' dreams."},
	{"fd3fb0a7ffb8af16603f3d3af98f8e1f", "You remind me of a TV show, but that's all right: I watch it anyway."},
	{"469b13a78ebf297ecda64d4723655154", "C is as portable as Stonehedge!!"},
	{"63eb3a2f466410104731c4b037600110", "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{"72c2ed7592debca1c90fc0100f931a2f", "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{"132f7619d33b523b1d9e5bd8e0928355", "How can you write a big system without C++?  -Paul Glick"},
}

func testEqual(t *testing.T, msg string, args ...interface{}) bool {
	t.Helper()
	if args[len(args)-2] != args[len(args)-1] {
		t.Errorf(msg, args...)
		return false
	}
	return true
}

func TestMbase64Encode(t *testing.T) {
	for _, p := range pairs {
		bb := &bytes.Buffer{}
		err := ydocr.Mbase64Encode([]byte(p.decoded), bb)
		if err != nil {
			t.Errorf("Encode(%q) = %q, err: %s ", p.decoded, p.encoded, err)
		}
		testEqual(t, "Encode(%q) = %q, want %q", p.decoded, bb.String(), p.encoded)
	}
}

func TestCountMd5(t *testing.T) {
	for _, p := range golden {
		s := ydocr.MCountMd5([]byte(p.in))
		if s != p.out {
			t.Fatalf("Sum function: md5(%s) = %s want %s", p.in, s, p.out)
		}
	}
}


func TestCountSign(t *testing.T) {
	img := "nasdfljliena123nlnafu2o3nAFaflfo2nl4nAf2134ljalsjf"

	sign, salt, err := ydocr.CountSign(appKey, appSecret, img)
	if err != nil {
		t.Fatal("count sign err: %s", err)
	}

	t.Logf("sign %s, salt %s", sign, salt)
}

func TestOcr_BasePost(t *testing.T) {
	o := ydocr.NewOcr(appKey, appSecret)

	file, _ := os.Open(filename)
	f, _ := ioutil.ReadAll(file)

	rsp, err := o.BasePost(f)
	testCommonLine(t, rsp, err)
}

func TestOcr_FileOcr(t *testing.T) {
	o := ydocr.NewOcr(appKey, appSecret)
	rsp, err := o.FileOcr(filename)
	testCommonLine(t, rsp, err)

}

func TestOcr_ImageOcr(t *testing.T) {
	o := ydocr.NewOcr(appKey, appSecret)
	i, err := fileToImage(filename)
	rsp, err := o.ImageOcr(&i)
	testCommonLine(t, rsp, err)
}

func testCommonLine(t *testing.T, rsp []byte, err error) {
	if err != nil{
		log.Fatal("base post err: %s", err)
	}
	l, err := ydocr.Response2Line(rsp)
	if err != nil{
		log.Fatal("response to line err: %s", err)
	}
	t.Logf("rsp: %+v", l)
}


func fileToImage(p string) (image.Image, error){
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open %s err: %s", p, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s err: %s", p, err)
	}
	return img, nil
}
