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

	NSString* myString = [myPasteboard  stringForType:bestType];

	if (bestType == (id)[NSNull null] || IsEmpty(myString)) {
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
)

type CBType struct {
	Type string
	Data string
}

var (
	old = ""
)

func ClipBoardStart() {
	for _ = range time.Tick(time.Second) {
		data := C.GetPaste()
		str := C.GoString(data)

		if len(str) > 0 {
			h256 := sha256.New()
			io.WriteString(h256, str)
			sha := fmt.Sprintf("%x", h256.Sum(nil))

			cb := CBType{
				Data: str,
				Type: "",
			}

			if sha != old {
				old = sha

				if valid.IsEmail(str) {
					cb.Type = "email"
				} else if valid.IsURL(str) {
					cb.Type = "url"
				} else if valid.IsJSON(str) {
					cb.Type = "json"
				} else if valid.IsIP(str) {
					if valid.IsIPv4(str) {
						cb.Type = "ipv4"
					} else if valid.IsIPv6(str) {
						cb.Type = "ipv6"
					}
				} else {
					cb.Type = "none"
				}

				OutputChan <- cb

			}
		}

	}
}
