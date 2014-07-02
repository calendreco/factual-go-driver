package factual_test

import (
	"log"
	"github.com/calendreco/factual"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Suite")
}

var _ = Describe("Factuals api", func(){

	instance := factual.New("h29FhoBkKNovV6DuBNoiOHjl3YSD6sj2ksRgthnD",
						    "Z9bjczytqVXfMxluYNgO8kAOgtLZJavHI8JHEo3M")

	Describe("The places api", func(){
		t := instance.Table("places")

		It("Should return individual entries", func(){
			c := t.Id("5a46e853-a617-4ce6-8bd9-de0daa3c76f4")
			f := factual.Place{}
			err := c.Iter().One(&f)
			Ω(err).Should(BeNil())
			Ω(f.Name).Should(Equal("Chipotle Mexican Grill"))
		})

		It("Should support searching", func(){
			c := t.Search("Chipotle")
			ps := []factual.Place{}
			err := c.Iter().All(&ps)
			Ω(err).Should(BeNil())
			for _, p := range ps{
				log.Println(p.Name)
				Ω(p.Name).Should(ContainSubstring("Chipotle"))
			}
		})

		It("Should support filtering", func(){
			c := t.Filter(factual.F{"region":"NY"})
			ps := []factual.Place{}
			err := c.Iter().All(&ps)
			Ω(err).Should(BeNil())
			for _, p := range ps{
				log.Println(p.Region)
				Ω(p.Region).Should(Equal("NY"))
			}
		})

	})

	FDescribe("The crosswalk api", func(){
		t := instance.Table("crosswalk")

		It("Should return for a specific id", func(){
			c := t.Id("5a46e853-a617-4ce6-8bd9-de0daa3c76f4")

			cs := []factual.Crosswalk{}
			err := c.Iter().All(&cs)
			Ω(err).Should(BeNil())

			for _, c := range cs{
				log.Println(c.FactualId)
				Ω(c.FactualId).Should(Equal("5a46e853-a617-4ce6-8bd9-de0daa3c76f4"))
			}
		})
	})

})