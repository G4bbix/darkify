package main

import (
  "fmt"
)

func main() {

  /* default bg color, will be the replacement for white
     based on this all other colors will be calculated */
  input_file_path, dark_bg_rgb := Parse_args()
  // fmt.Println(dark_bg_rgb)
  inputLines := Read_input_file(input_file_path)

  light_css_vars := Extract_css_vars(inputLines)

  var light_css_vars_rgb = make(map[string][3]uint8)
  for key, value := range light_css_vars {
    light_css_vars_rgb[key] = Hex_str_to_int_arr(value)
  }
  light_bg_rgb := [3]uint8 {255, 255, 255}

  dark_css_vars_rgb := Offset_algo(light_css_vars_rgb, dark_bg_rgb, light_bg_rgb, 100)

  for key, value := range dark_css_vars_rgb {
    fmt.Println("--" + key + ": " + Int_arr_to_hex_str(value))
  }
}
