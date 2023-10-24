package main

import "fmt"

func main() {

	/* default bg color, will be the replacement for white
	   based on this all other colors will be calculated */
	input_file_path, dark_bg_rgb := Parse_args()
	light_css_vars := Read_input_file(input_file_path)
	var light_bg_rgb = [3]uint8{255, 255, 255}

	mid_point_sum := Calc_midpoint_sum(light_bg_rgb, dark_bg_rgb)
	// fmt.Println(mid_point_sum)
	var dark_css_vars_rgb = make(map[string][3]uint8)
	for key, value := range light_css_vars {
		dark_css_vars_rgb[key] = Darkify(value, mid_point_sum)
	}
	fmt.Println("@media (prefers-color-schmeme: dark) {")
	for key, value := range dark_css_vars_rgb {
		fmt.Printf("\t--%s: %s\n", key, Int_arr_to_hex_str(value))
	}
	fmt.Println("}")
}
