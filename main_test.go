package main

import (
	"reflect"
	"testing"
)

func TestComputeRevertFilter(t *testing.T) {
	revertFilter := computeRevertFilter(1)
	if revertFilter[0] != 0 {
		t.Errorf("revertFilter(1) shoulf return 0, returns %v", revertFilter[0])
	}
	revertFilter = computeRevertFilter(7)
	want := []uint32{6, 5, 4, 3, 2, 1, 0}
	if !reflect.DeepEqual(revertFilter, want) {
		t.Errorf("want : %v ; obtain : %v\n", want, revertFilter)
	}
}

func TestComputeZigzagFilter(t *testing.T) {
	/*
		matrix := [][]uint32{
			{2, 3, 8, 9, 14},
			{1, 4, 7, 10, 13},
			{0, 5, 6, 11, 12}}
	*/
	want := []uint32{10, 5, 0, 1, 6, 11, 12, 7, 2, 3, 8, 13, 14, 9, 4}
	filter := computeZigzagFilter(5, 3)
	if !reflect.DeepEqual(filter, want) {
		t.Errorf("want : %v ; obtain : %v\n", want, filter)
	}
}

func Test_matrixToLinear(t *testing.T) {
	matrix := [][]uint32{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		{26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38},
		{39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51},
		{52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
		{65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77},
		{78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}}
	want := []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	line := matrixToLinear(matrix)
	if !reflect.DeepEqual(line, want) {
		t.Errorf("want : %v ; obtain : %v\n", want, line)
	}
}

func Test_applyFilter(t *testing.T) {
	line := []uint32{10, 11, 12, 13, 14}
	filter := []uint32{1, 2, 0, 4, 3}

	want := []uint32{11, 12, 10, 14, 13}
	if !reflect.DeepEqual(applyFilter(line, filter), want) {
		t.Errorf("want : %v ; obtain : %v\n", want, filter)
	}
}

func Test_computeRollXFilter(t *testing.T) {
	matrix := [][]uint32{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		{26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38},
		{39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51},
		{52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
		{65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77},
		{78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}}

	line := matrixToLinear(matrix)
	wantx := [][]uint32{
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0},
		{14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 13},
		{27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 26},
		{40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 39},
		{53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 52},
		{66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 65},
		{79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 78}}
	want := matrixToLinear(wantx)
	result := applyFilter(line, computeRollXFilter(13, 7))
	if !reflect.DeepEqual(result, want) {
		t.Errorf("\n  line : %v ; \n  want : %v ; \nobtain : %v\nwith    :%v\n", line, want, result, computeRollXFilter(13, 7))
	}
}
