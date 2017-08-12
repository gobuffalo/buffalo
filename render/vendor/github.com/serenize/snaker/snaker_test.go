package snaker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Snaker", func() {
	Describe("CamelToSnake test", func() {
		It("should return an empty string on an empty input", func() {
			Expect(CamelToSnake("")).To(Equal(""))
		})

		It("should work with one word", func() {
			Expect(CamelToSnake("One")).To(Equal("one"))
		})

		It("should return an uppercase string as seperate words", func() {
			Expect(CamelToSnake("ONE")).To(Equal("o_n_e"))
		})

		It("should return ID as lowercase", func() {
			Expect(CamelToSnake("ID")).To(Equal("id"))
		})

		It("should work with a single lowercase character", func() {
			Expect(CamelToSnake("i")).To(Equal("i"))
		})

		It("should work with a single uppcase character", func() {
			Expect(CamelToSnake("I")).To(Equal("i"))
		})

		It("should return a long text as expected", func() {
			Expect(CamelToSnake("ThisHasToBeConvertedCorrectlyID")).To(
				Equal("this_has_to_be_converted_correctly_id"))
		})

		It("should return the text as expected if the initialism is in the middle", func() {
			Expect(CamelToSnake("ThisIDIsFine")).To(Equal("this_id_is_fine"))
		})

		It("should work with long initialism", func() {
			Expect(CamelToSnake("ThisHTTPSConnection")).To(Equal("this_https_connection"))
		})

		It("should work with multi initialisms", func() {
			Expect(CamelToSnake("HelloHTTPSConnectionID")).To(Equal("hello_https_connection_id"))
		})

		It("sould work with concat initialisms", func() {
			Expect(CamelToSnake("HTTPSID")).To(Equal("https_id"))
		})
	})

	Describe("SnakeToCamel test", func() {
		It("should return an empty string on an empty input", func() {
			Expect(SnakeToCamel("")).To(Equal(""))
		})

		It("should not blow up on trailing _", func() {
			Expect(SnakeToCamel("potato_")).To(Equal("Potato"))
		})

		It("should return a snaked text as camel case", func() {
			Expect(SnakeToCamel("this_has_to_be_uppercased")).To(
				Equal("ThisHasToBeUppercased"))
		})

		It("should return a snaked text as camel case, except the word ID", func() {
			Expect(SnakeToCamel("this_is_an_id")).To(Equal("ThisIsAnID"))
		})

		It("should return 'id' not as uppercase", func() {
			Expect(SnakeToCamel("this_is_an_identifier")).To(Equal("ThisIsAnIdentifier"))
		})

		It("should simply work with id", func() {
			Expect(SnakeToCamel("id")).To(Equal("ID"))
		})
	})

	Describe("SnakeToCamelLower test", func() {
		It("should return an empty string on an empty input", func() {
			Ω(SnakeToCamelLower("")).To(Equal(""))
		})

		It("should not blow up on trailing _", func() {
			Ω(SnakeToCamelLower("potato_")).To(Equal("potato"))
		})

		It("should return a snaked text as camel case", func() {
			Ω(SnakeToCamelLower("this_has_to_be_uppercased")).To(
				Equal("thisHasToBeUppercased"))
		})

		It("should return a snaked text as camel case, except the word ID", func() {
			Ω(SnakeToCamelLower("this_is_an_id")).To(Equal("thisIsAnID"))
		})

		It("should return 'id' not as uppercase", func() {
			Ω(SnakeToCamelLower("this_is_an_identifier")).To(Equal("thisIsAnIdentifier"))
		})

		It("should simply work with id", func() {
			Ω(SnakeToCamelLower("id")).To(Equal("id"))
		})

		It("should simply work with leading id", func() {
			Ω(SnakeToCamelLower("id_me_please")).To(Equal("idMePlease"))
		})
	})
})
