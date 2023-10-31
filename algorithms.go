package main

import (
	"math"
)

func get_amount_zeros(color [3]uint8) uint8 {

	var amount_zeros uint8
	for _, segment := range color {
		if segment == 0 {
			amount_zeros++
		}
	}
	return amount_zeros
}

func calc_remaining_unused_space(color [3]uint8) uint16 {
	var remaining_unused_space uint16
	for _, segment := range color {
		if segment != 255 {
			remaining_unused_space += uint16(255 - segment)
		}
	}
	return remaining_unused_space
}

func calc_leftover_dist_two_zeros(color [3]uint8) [3]float64 {
	var dist [3]float64
	for i, segment := range color {
		if segment != 255 {
			dist[i] = 0.5
		} else {
			dist[i] = 0
		}
	}
	return dist
}

// With one zero e.G. [255, 0, 128] we use 255 - segment as the dist template sum -> (255-255) + (255-0) + (255-128)
func calc_dist_template_sum(color [3]uint8) uint16 {
	var dist_template uint16
	for _, segment := range color {
		dist_template += (255 - uint16(segment))
	}
	return dist_template
}

func calc_leftover_dist_one_zero(color [3]uint8, remaining_unused_space uint16) [3]float64 {
	var dist [3]float64
	segment_to_flip := 4
	for i, segment := range color {
		if segment != 255 {
			if segment_to_flip == 4 {
				dist[i] = float64(255-segment) / float64(remaining_unused_space)
				segment_to_flip = i
			} else {
				dist[i] = dist[segment_to_flip]
				dist[segment_to_flip] = float64(255-segment) / float64(remaining_unused_space)
			}
		} else {
			dist[i] = 0
		}
	}
	return dist
}

func calc_leftover_dist_no_zeros(color [3]uint8, remaining_unused_space uint16) [3]float64 {
	var dist [3]float64
	for i, segment := range color {
		if segment != 255 {
			dist[i] = float64(segment) / float64(remaining_unused_space)
		} else {
			dist[i] = 0
		}
	}
	return dist
}

// Get the percieved birghtness of the sum of the point between the light bg and dark bg to ITU BT.601
func Calc_midpoint_sum(light_bg_rgb [3]uint8, dark_bg_rgb [3]uint8) float32 {
	return (((float32(light_bg_rgb[0]) + float32(dark_bg_rgb[0])) * 0.299) +
		((float32(light_bg_rgb[1]) + float32(dark_bg_rgb[1])) * 0.587) +
		((float32(light_bg_rgb[2]) + float32(dark_bg_rgb[2])) * 0.114)) * 3 / 2
}

// Calculate percieved brightness according to ITU BT.601
func Calc_percieved_brightness(rgb [3]uint8) uint8 {
	return uint8(float64(rgb[0])*0.299 + float64(rgb[1])*0.587 + float64(rgb[2])*0.114)
}

func Darkify(light_rgb [3]uint8, mid_point_sum float32) [3]uint8 {
	var sum_light_rgb uint16 = uint16(light_rgb[0]) + uint16(light_rgb[1]) + uint16(light_rgb[2])
	var sum_dark_rgb uint16 = uint16(calc_dark_rgb_sum(sum_light_rgb, mid_point_sum))
	distribution := calc_distribution(light_rgb, sum_light_rgb)
	dark_rgb, leftover := initCalc(sum_dark_rgb, light_rgb, distribution)
	if leftover > 0 {
		amount_zeros := get_amount_zeros(dark_rgb)

		remaining_unused_space := calc_remaining_unused_space(dark_rgb)

		var dist [3]float64
		if amount_zeros == 1 {
			dist_template := calc_dist_template_sum(dark_rgb)
			dist = calc_leftover_dist_one_zero(dark_rgb, dist_template)
		} else if amount_zeros == 2 {
			dist = calc_leftover_dist_two_zeros(dark_rgb)
		} else {
			dist = calc_leftover_dist_no_zeros(dark_rgb, remaining_unused_space)
		}
		distribute_leftovers(&dark_rgb, leftover, dist)
	}
	return dark_rgb
}

func calc_dark_rgb_sum(sum_light_rgb uint16, mid_point_sum float32) int16 {
	var offset int16 = int16(sum_light_rgb) - int16(mid_point_sum*2)
	// if offset is negative, invert the algabreic sign
	if offset < 0 {
		return offset / -1
	} else {
		return offset
	}
}

func calc_distribution(rgb [3]uint8, sum uint16) [3]float64 {

	if rgb[0] == 0 && rgb[1] == 0 && rgb[2] == 0 {
		var distribution = [3]float64{0.33333334, 0.33333334, 0.33333334}
		return distribution
	}

	var distribution [3]float64
	for i := 0; i < 3; i++ {
		distribution[i] = float64(rgb[i]) / float64(sum)
	}
	return distribution
}

func initCalc(dark_rgb_sum uint16, light_rgb [3]uint8, dist [3]float64) ([3]uint8, uint16) {
	var dark_rgb [3]uint8
	var leftover uint16 = 0
	var val uint16

	for i := range light_rgb {
		val = uint16(math.Round(float64(dist[i] * float64(dark_rgb_sum))))
		if val > 255 {
			leftover += val - 255
			dark_rgb[i] = 255
		} else {
			dark_rgb[i] = uint8(val)
		}
	}
	return dark_rgb, leftover
}

func distribute_leftovers(color *[3]uint8, leftover uint16, dist [3]float64) {

	// Calculate new distribution
	var val uint16
	var leftover_distributed uint16
	var color_with_leftover [3]uint8

	for i := range *color {
		if dist[i] != 0 {
			leftover_to_dist := uint16(float64(leftover) * dist[i])
			val = leftover_to_dist + uint16(color[i])
			if val > 255 {
				leftover_distributed += val - uint16(color[i])
				color[i] = 255
			} else {
				color[i] = uint8(val)
				leftover_distributed += leftover_to_dist
			}
		}
	}

	var remaining_leftover uint16 = leftover - leftover_distributed
	if remaining_leftover > 0 {
		for i, segment := range *color {
			if segment != 255 {
				val = uint16(color[i]) + remaining_leftover
				if val >= 255 {
					color_with_leftover[i] = 255
				} else {
					color_with_leftover[i] = uint8(val)
				}
			}
		}
	}
}
