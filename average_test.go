package main

import (
    "testing"
)

func TestAverage(t *testing.T) {
    tests := []struct {
        name     string
        numbers  []float64
        expected float64
    }{
        {"average of positive numbers", []float64{1.0, 2.0, 3.0, 4.0}, 2.5},
        {"average of mixed numbers", []float64{-1.0, -2.0, 2.0, 4.0}, 0.75},
        {"empty slice", []float64{}, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Average(tt.numbers)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}