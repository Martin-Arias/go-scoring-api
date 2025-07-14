package utils_test

import (
	"github.com/Martin-Arias/go-scoring-api/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Statistics calculation", func() {

	It("should return mean, mode and median for an array of odd length", func() {
		points := []int{200, 100, 20, 20, 80}
		mean, median, mode := utils.CalculateStatistics(points)
		Expect(mean).To(Equal(float64(84)))
		Expect(median).To(Equal(float64(80)))
		Expect(mode).To(Equal([]int{20}))
	})

	It("should return zero for all values", func() {
		points := []int{}
		mean, median, mode := utils.CalculateStatistics(points)
		Expect(mean).To(Equal(float64(0)))
		Expect(median).To(Equal(float64(0)))
		Expect(mode).To(Equal([]int{}))
	})

	It("should return mean, mode and median for an array of odd length", func() {
		points := []int{200, 100, 20, 30, 20, 80}
		mean, median, mode := utils.CalculateStatistics(points)
		Expect(mean).To(Equal(float64(75)))
		Expect(median).To(Equal(float64(55)))
		Expect(mode).To(Equal([]int{20}))
	})

})
