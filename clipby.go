package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

static inline BOOL IsEmpty(id thing) {
	return thing == nil
	|| ([thing respondsToSelector:@selector(length)]
	&& [(NSData *)thing length] == 0)
	|| ([thing respondsToSelector:@selector(count)]
	&& [(NSArray *)thing count] == 0);
}

const char* GetPaste() {

	NSArray *supportedTypes = [NSArray arrayWithObjects: NSPasteboardTypeRTF, NSPasteboardTypeRTFD,
	NSPasteboardTypeRuler, NSPasteboardTypeMultipleTextSelection, NSRTFPboardType,
	NSRTFDPboardType, NSPasteboardTypeTabularText, NSStringPboardType, nil];

	NSPasteboard*  myPasteboard  = [NSPasteboard pasteboardWithName:NSGeneralPboard];

	NSString *bestType = [myPasteboard availableTypeFromArray:supportedTypes];

	// NSLog( @"%@", bestType );

	NSString* myString = [myPasteboard  stringForType:bestType];

	if (bestType == (id)[NSNull null] || IsEmpty(myString)) {
		NSLog( @"%@", bestType );
		myString = @"";
	}

    return [myString UTF8String];
}
*/
import "C"
import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/mitchellh/go-mruby"
)

var (
	old = ""
)

const ()

func main() {
	mrb := mruby.NewMrb()
	defer mrb.Close()

	// Our custom function we'll expose to Ruby. The first return
	// value is what to return from the func and the second is an
	// exception to raise (if any).
	addFunc := func(m *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		args := m.GetArgs()
		return mruby.Int(args[0].Fixnum() + args[1].Fixnum()), nil
	}

	// Lets define a custom class and a class method we can call.
	class := mrb.DefineClass("Example", nil)
	class.DefineClassMethod("add", addFunc, mruby.ArgsReq(2))

	// Let's call it and inspect the result
	result, err := mrb.LoadString(`Example.add(12, 30)`)
	if err != nil {
		panic(err.Error())
	}

	// This will output "Result: 42"
	fmt.Printf("Result: %s\n", result.String())
}

func Run() {

	for _ = range time.Tick(time.Second) {
		data := C.GetPaste()
		str := C.GoString(data)

		if len(str) > 0 {
			h256 := sha256.New()
			io.WriteString(h256, str)
			sha := fmt.Sprintf("%x", h256.Sum(nil))

			if sha != old {
				old = sha

				if valid.IsEmail(str) {
					fmt.Println("GOT EMAIL: ", str)
				} else if valid.IsURL(str) {
					fmt.Println("GOT URL: ", str)
				} else if valid.IsJSON(str) {
					fmt.Println("GOT JSON: ", str)
				} else if valid.IsIP(str) {
					fmt.Println("GOT IP: ", str)

					if valid.IsIPv4(str) {
						fmt.Println("GOT IPv4: ", str)
					} else if valid.IsIPv6(str) {
						fmt.Println("GOT IPv6: ", str)
					}
				}

			} else {
				// fmt.Println(sha, "\n")
			}
		}

	}
}
