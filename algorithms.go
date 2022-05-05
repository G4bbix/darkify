package main

// contains the definiton of all color manipulating algorythms

func Offset_algo(light_vars map[string][3]uint8, dark_bg [3]uint8, light_bg [3]uint8,
  use_remaining uint8) map[string][3]uint8 {

  var dark_css_vars = make(map[string][3]uint8)
  var leftover uint16
  for key, value := range light_vars {
    if value == light_bg {
      dark_css_vars[key] = dark_bg
      continue
    }
    var sum_dark_color float32 = 0
    for i := 0; i < 3; i++ {
      if value[i] == 0 {
        continue
      } else {
        sum_dark_color += float32(light_bg[i]) * float32(dark_bg[i]) / float32(value[i])
      }
    }
    // Workaround for edgecase of very dark colors.
    // This should be replaced if modification other than darkening are wanted
    if sum_dark_color == 0 || sum_dark_color > 765{
      sum_dark_color = 765
    }
    var sum_light_color float32 = float32(uint16(value[0]) + uint16(value[1]) + uint16(value[2]))

    // calculate color_dark and save leftover
    var color_dark [3]uint8
    for i, octett := range value {
      color_dark[i] = calc_dark_val(sum_light_color, sum_dark_color, float32(octett))
    }

    // Reduce leftover by using the percent specified
    leftover = uint16((sum_dark_color - float32(uint16(color_dark[0]) + uint16(color_dark[1]) +
      uint16(color_dark[2]))) / (100 * float32(use_remaining)))

    if leftover > 0 {
      dark_css_vars[key] = distribute_leftover(leftover, value)
    } else {
      dark_css_vars[key] = color_dark
    }
  }
  return dark_css_vars
}

// calculates the new value of an RGB octett
func calc_dark_val(sum_light_color float32, sum_dark_color float32,
  light_value float32) uint8 {

  var share float32 = light_value / sum_light_color
  var dark_value uint16 = uint16(share * sum_dark_color)

  if dark_value > 255 {
    return 255
  } else {
    return uint8(dark_value)
  }
}

func distribute_leftover(leftover uint16, color_dark [3]uint8) [3]uint8 {
  var total uint16
  for _, value := range color_dark {
    if value < 255 {
      total += uint16(value)
    }
  }
  var color_dark_changed [3]uint8
  var remaining_after uint8 = 0
  if color_dark[0] < 255 {
    var red_result [2]uint8 = calclate_and_add_additon(leftover, color_dark[0])
    color_dark_changed[0] = red_result[1]
    remaining_after += red_result[0]
  }
  if color_dark[1] < 255 {
    var green_result [2]uint8 = calclate_and_add_additon(leftover, color_dark[1])
    color_dark_changed[1] = green_result[1]
    remaining_after += green_result[0]
  }
  if color_dark[2] < 255 {
    var blue_result [2]uint8 = calclate_and_add_additon(leftover, color_dark[2])
    color_dark_changed[2] = blue_result[1]
    remaining_after += blue_result[0]
  }

  if remaining_after > 0 {
    for i := 0; i > 3; i++ {
      if color_dark_changed[i] != 255 {
        if color_dark_changed[i] + remaining_after >= 255 {
          color_dark_changed[i] = 255
        } else {
          color_dark_changed[i] += remaining_after
        }
      }
    }
  }
  return color_dark_changed
}

func calclate_and_add_additon(total uint16, color uint8) [2]uint8 {
  var divisor uint16
  if color == 0 {
    divisor = 1
  } else {
    divisor = uint16(color)
  }
  color_share := total / divisor
  addition := color_share * total
  var color_remaining [2]uint8
  if addition > 255 {
    color_remaining[0] = uint8(addition - 255)
    color_remaining[1] = 255
  } else {
    color_remaining[0] = 0
    color_remaining[1] = color
  }
  return color_remaining
}
