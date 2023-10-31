package main

import (
	"fmt"
	"strings"
)

func format_css_var(name string, value [3]uint8, format string, indent_level int) string {
	var value_formatted string
	if format == "hex" {
		value_formatted = Int_arr_to_hex_str(value)
	} else if format == "rgb" {
		value_formatted = fmt.Sprintf("rgb(%d,%d,%d)", value[0], value[1], value[2])
	}
	return fmt.Sprintf("%s--%s: %s;\n", strings.Repeat("\t", indent_level), name, value_formatted)
}

func main() {

	/* default bg color, will be the replacement for white
	   based on this all other colors will be calculated */
	input_file_path, dark_bg_rgb, output_format := Parse_args()
	light_css_vars := Read_input_file(input_file_path)
	var light_bg_rgb = [3]uint8{255, 255, 255}

	mid_point_sum := Calc_midpoint_sum(light_bg_rgb, dark_bg_rgb)
	// fmt.Println(mid_point_sum)
	var dark_css_vars_rgb = make(map[string][3]uint8)
	fmt.Println(":root {")
	for key, value := range light_css_vars {
		fmt.Printf("%s", format_css_var(key, value, output_format, 1))
		dark_css_vars_rgb[key] = Darkify(value, mid_point_sum)
	}
	fmt.Println("}")

	fmt.Printf("@media (prefers-color-schmeme: dark) {\n\t:root {\n")
	for key, value := range dark_css_vars_rgb {
		fmt.Printf("%s", format_css_var(key, value, output_format, 2))
	}
	fmt.Println("\t}\n}")
}
