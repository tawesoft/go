package dialog // import "tawesoft.co.uk/go/dialog"

import (
	"testing"
)

func TestWraps(t *testing.T) {
    var tests = [][3]interface{}{
        {"",  "",   1},
        {"a", "a", -1},
        {"a", "a",  0},
        {"a", "a",  1},
        {"a", "a",  2},
        {"  a  ", "a",  1},
        {"hello\nworld", "hello\nworld",  2},
        {"hello world", "hello\nworld",  2},
        {"hello world", "hello\nworld",  5},
        {"hello world", "hello\nworld", 10},
        {"hello world", "hello world",  11},
        {"hello world", "hello world",  12},
        {"hello\nworld", "hello world",  12},
        {"hello      world", "hello world", 12},
        {"    hello    world   ", "hello world", 12},
        {"a b c d e f g h i", "a\nb\nc\nd\ne\nf\ng\nh\ni", -1},
        {"a b c d e f g h i", "a\nb\nc\nd\ne\nf\ng\nh\ni", 0},
        {"a b c d e f g h i", "a\nb\nc\nd\ne\nf\ng\nh\ni", 1},
        {"a b c d e f g h i", "a b c\nd e f\ng h i", 5},
        {"a b c d e f g h i", "a b c\nd e f\ng h i", 6},
        {"a b c d e f g h i", "a b c d e\nf g h i",  9},
        {"a b c d e f g h i", "a b c d e\nf g h i", 10},
        {`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce a
        tortor sagittis, elementum velit id, scelerisque erat. Sed mollis odio
        molestie dui venenatis condimentum. Donec massa ligula, auctor rutrum
        interdum a, faucibus sed sapien. Vivamus neque massa, porttitor vel
        nulla eu, gravida egestas massa. Aliquam interdum pellentesque elit.
        Quisque vestibulum, libero condimentum venenatis commodo, erat lectus
        convallis libero, at pellentesque nibh enim vel risus. Duis elit mi,
        lacinia ut ex vitae, ullamcorper tempus ex. Lorem ipsum dolor sit amet,
        consectetur adipiscing elit. Fusce eu elit molestie, tempor nulla
        vehicula, tempor nulla. Maecenas pellentesque, lectus non accumsan
        pharetra, neque justo dignissim dolor, sit amet luctus mi leo ut dui.`,

        `Lorem ipsum dolor sit amet,
consectetur adipiscing elit.
Fusce a tortor sagittis,
elementum velit id,
scelerisque erat. Sed mollis
odio molestie dui venenatis
condimentum. Donec massa
ligula, auctor rutrum interdum
a, faucibus sed sapien.
Vivamus neque massa, porttitor
vel nulla eu, gravida egestas
massa. Aliquam interdum
pellentesque elit. Quisque
vestibulum, libero condimentum
venenatis commodo, erat lectus
convallis libero, at
pellentesque nibh enim vel
risus. Duis elit mi, lacinia
ut ex vitae, ullamcorper
tempus ex. Lorem ipsum dolor
sit amet, consectetur
adipiscing elit. Fusce eu elit
molestie, tempor nulla
vehicula, tempor nulla.
Maecenas pellentesque, lectus
non accumsan pharetra, neque
justo dignissim dolor, sit
amet luctus mi leo ut dui.`,
        30},
    }

    for index, test := range tests {
        var original, expected, length = test[0], test[1], test[2]
        var result = wrap(original.(string), length.(int))
        if result != expected.(string) {
            t.Errorf("Test %d failed: wrap(\"%s\", %d), got \"%s\" but wanted \"%s\"\n",
                index, original.(string), length.(int), result, expected.(string))
        }
    }
}
