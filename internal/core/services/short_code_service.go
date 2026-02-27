package services

import (
	"URLShortner/internal/core/domain"
	"URLShortner/pkg"
	"context"
	crand "crypto/rand"
	"errors"
	"math/big"
	"math/rand/v2"
	"strings"
)

type ShortCodeService struct {
	generator   *CodeGenerator
	codeChannel chan string
	repository  domain.UrlRepository
}

func NewCodeGeneratorService(urlRepo domain.UrlRepository, config pkg.ShortCodeConfig) *ShortCodeService {
	codeGen := ShortCodeService{
		generator:   NewCodeGenerator(config),
		codeChannel: make(chan string),
		repository:  urlRepo,
	}
	go codeGen.produceShortCode()
	return &codeGen
}

type CodeGenerator struct {
	config pkg.ShortCodeConfig
}

func NewCodeGenerator(config pkg.ShortCodeConfig) *CodeGenerator {
	return &CodeGenerator{
		config: config,
	}
}

func randomInt(max int) int {
	n, err := crand.Int(crand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return rand.IntN(max)
	}
	return int(n.Int64())
}

func (g *CodeGenerator) randomPart(randomLength int) string {
	alphabetSwitcher := rand.IntN(2) == 0
	alphabetCounter := 0

	var builder strings.Builder
	builder.Grow(randomLength)
	for i := 0; i < randomLength; i++ {
		if alphabetSwitcher && alphabetCounter < 2 {
			builder.WriteByte(g.config.Chars[randomInt(len(g.config.Chars))])
			alphabetCounter++
		} else {
			builder.WriteByte(g.config.Digits[randomInt(len(g.config.Digits))])
			alphabetCounter = 0
		}
		alphabetSwitcher = rand.IntN(2) == 0
	}

	return builder.String()
}

func (g *CodeGenerator) Generate() []string {
	randomString := g.randomPart(g.config.MaxLength)
	randomRunes := []rune(randomString)

	result := make([]string, 0, g.config.MaxLength-g.config.MinLength+1)

	for length := g.config.MinLength; length <= len(randomRunes); length++ {
		result = append(result, string(randomRunes[:length]))
	}
	return result
}

func (c *ShortCodeService) produceShortCode() {
	for {
		candidates := c.generator.Generate()
		for _, candidate := range candidates {
			_, err := c.repository.GetUrlById(context.Background(), candidate)
			if errors.Is(err, domain.ErrUrlNotFound) {
				c.codeChannel <- candidate
				break
			}
		}
	}
}

func (c *ShortCodeService) GetShortCode() string {
	return <-c.codeChannel
}
