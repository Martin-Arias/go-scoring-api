package utils_test

import (
	"github.com/Martin-Arias/go-scoring-api/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Statistics calculation", func() {

	It("should return mean, mode and median", func() {

		points := []int{200, 100, 20, 20, 80}

		mean, median, mode := utils.CalculateStatistics(points)

		Expect(mean).To(Equal(float64(84)))
		Expect(median).To(Equal(float64(80)))
		Expect(mode).To(Equal([]int{20}))
	})

})
