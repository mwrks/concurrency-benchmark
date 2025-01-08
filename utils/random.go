package utils

import (
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

var paragraph = "Lorem ipsum dolor sit amet consectetur adipiscing elit Nullam vehicula ex id quam tincidunt ac varius justo cursus Proin ac efficitur risus quis dapibus tortor Praesent sit amet vehicula lorem vel pharetra mauris Aenean congue felis a sapien ultricies hendrerit Curabitur in sem vitae mi sagittis bibendum in nec elit Cras vel nisl vel risus dictum tincidunt vel id libero Sed aliquet dolor eget libero aliquet vel aliquet sem consequat Vivamus auctor justo in urna gravida faucibus Fusce luctus purus vel pharetra efficitur velit sapien tincidunt sapien eget vulputate quam ligula sed turpis Sed viverra hendrerit purus id posuere Ut quis finibus magna Aliquam sodales odio sed consequat maximus justo justo egestas lectus non commodo nisi sapien non ipsum Nullam non magna ut ligula accumsan fermentum Integer pellentesque velit eu orci aliquet id pharetra erat mollis Ut volutpat ligula nec ipsum fermentum sed interdum metus vehicula Suspendisse ac sapien at justo pharetra auctor in sed nisi Morbi molestie eros vel mauris tempor sodales Maecenas scelerisque erat id sapien aliquet vehicula Vestibulum scelerisque nisi sed rutrum scelerisque nisl enim aliquet dolor vel tincidunt sapien nulla et arcu Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas Ut at lectus non magna convallis blandit eget id tortor Curabitur gravida justo at lacinia dictum Duis malesuada lacinia quam nec cursus neque facilisis nec Donec efficitur suscipit tellus Quisque scelerisque orci et arcu vestibulum fermentum Fusce eget nulla nisl Cras vehicula sagittis tellus sit amet eleifend Praesent tincidunt sem ac tortor finibus quis mollis purus tincidunt Aenean tincidunt nunc vel tincidunt venenatis Sed vitae lectus id dolor dictum vehicula id nec sapien Nulla ac nunc nec enim interdum dictum in a libero Suspendisse potenti Praesent eget lacus nec sapien malesuada gravida in ac justo Duis fringilla justo et augue venenatis luctus Pellentesque consectetur ipsum quis velit bibendum non posuere nisi ultricies"
var rnd = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
var names = strings.Fields(paragraph)

func GetRandomName() string {
	return names[rnd.Intn(len(names))]
}
