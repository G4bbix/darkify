package main

import (
	"fmt"
)

func main() {

	input_file_path, dark_bg_rgb, output_format, mode := Parse_args()

	light_css_vars := Read_input_file(input_file_path)
	light_bg_rgb := Rgb{name: "", color: [3]uint8{255, 255, 255}, alpha: 255}

	mid_point_sum := Calc_midpoint_sum(light_bg_rgb.color, dark_bg_rgb.color)
	var dark_rgb_color [3]uint8

	fmt.Printf("@media (prefers-color-schmeme: dark) {\n\t:root {\n")

	for _, value := range light_css_vars {

		sum_dark_rgb, distribution := Darkify_prep(value.color, mid_point_sum)

		if mode == "auto" {
			dark_rgb_color = Darkify_relative(sum_dark_rgb, distribution, mid_point_sum, true)
			dark_rgb_color = Darkify_linear(sum_dark_rgb, dark_rgb_color)
		} else if mode == "linear" {
			dark_rgb_color = Darkify_linear(sum_dark_rgb, value.color)
		} else if mode == "relative" {
			dark_rgb_color = Darkify_relative(sum_dark_rgb, distribution, mid_point_sum, false)
		}
		dark_rgb := Rgb{color: dark_rgb_color, alpha: value.alpha}
		fmt.Printf("%s", Format_css_var(value.name, dark_rgb, output_format, 1))
	}
	fmt.Println("}")
}
