package services

import (
	"URLShortner/internal/adapters/repository"
	"URLShortner/internal/infrastructure/persistence/ent/enttest"
	"URLShortner/pkg"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/require"

	"testing"
)

func TestShortCodeService_produceShortCode(t *testing.T) {
	entClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer entClient.Close()

	config, err := pkg.LoadConfig("../../../configs/test.yaml")
	require.NoError(t, err)
	urlRep := repository.NewPgUrlRepository(entClient)

	codeGenerator := NewCodeGeneratorService(urlRep, config.ShortCode)
	for i := 0; i < 100; i++ {
		select {
		case val := <-codeGenerator.codeChannel:
			strLen := len(val)
			if strLen < config.ShortCode.MinLength || strLen > config.ShortCode.MaxLength {
				t.Errorf("ShortCde %s length %d is outside allowed range [%d, %d]",
					val, strLen, config.ShortCode.MinLength, config.ShortCode.MaxLength)
			}
		case <-time.After(1 * time.Second):
			t.Error("Timeout: no data written to channel")
		}
	}
}
