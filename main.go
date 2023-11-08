package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Rgb struct {
	color [3]uint8
	alpha uint8
}

func round_float(value float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

func format_css_var(name string, value Rgb, format string, indent_level int) string {
	var value_formatted string
	if format == "hex" {
		value_formatted = Int_arr_to_hex_str(value.color)
		if value.alpha != 255 {
			value_formatted = fmt.Sprintf("%s%s", value_formatted, fmt.Sprintf("%x", value.alpha))
		}
	} else if format == "rgb" {
		if value.alpha == 255 {
			value_formatted = fmt.Sprintf("rgb(%d, %d, %d)", value.color[0], value.color[1], value.color[2])
		} else {
			alpha := strconv.FormatFloat(round_float(float64(value.alpha)/255, 2), 'f', -1, 32)
			value_formatted = fmt.Sprintf("rgba(%d, %d, %d, %s)", value.color[0], value.color[1], value.color[2], alpha)
		}
	}
	return fmt.Sprintf("%s--%s: %s;\n", strings.Repeat("\t", indent_level), name, value_formatted)
}

func main() {

	input_file_path, dark_bg_rgb, output_format := Parse_args()

	light_css_vars := Read_input_file(input_file_path)
	light_bg_rgb := Rgb{color: [3]uint8{255, 255, 255}, alpha: 255}

	mid_point_sum := Calc_midpoint_sum(light_bg_rgb.color, dark_bg_rgb.color)
	// fmt.Println(mid_point_sum)
	var dark_css_vars_rgb = make(map[string]Rgb)
	fmt.Println(":root {")
	for key, value := range light_css_vars {
		fmt.Printf("%s", format_css_var(key, value, output_format, 1))
		dark_css_vars_rgb[key] = Rgb{color: Darkify(value.color, mid_point_sum), alpha: value.alpha}
	}
	fmt.Println("}")

	fmt.Printf("@media (prefers-color-schmeme: dark) {\n\t:root {\n")
	for key, value := range dark_css_vars_rgb {
		fmt.Printf("%s", format_css_var(key, value, output_format, 2))
	}
	fmt.Println("\t}\n}")
}
