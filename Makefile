include .env.local
export

mocks:
	mockery --dir internal/ports --name IAlbumRepository --output internal/adapters/repository/mocks \
		--filename album.go --structname AlbumRepository
	mockery --dir internal/ports --name ICommentRepository --output internal/adapters/repository/mocks \
		--filename comment.go --structname CommentRepository
	mockery --dir internal/ports --name IGenreRepository --output internal/adapters/repository/mocks \
		--filename genre.go --structname GenreRepository
	mockery --dir internal/ports --name IMusicianRepository --output internal/adapters/repository/mocks \
		--filename musician.go --structname MusicianRepository
	mockery --dir internal/ports --name IOrderRepository --output internal/adapters/repository/mocks \
		--filename order.go --structname OrderRepository
	mockery --dir internal/ports --name ISubscriptionRepository --output internal/adapters/repository/mocks \
		--filename subscription.go --structname SubscriptionRepository
	mockery --dir internal/ports --name IUserRepository --output internal/adapters/repository/mocks \
        --filename user.go --structname UserRepository
	mockery --dir internal/ports --name IStatRepository --output internal/adapters/repository/mocks \
		--filename stat.go --structname StatRepository
	mockery --dir internal/ports --name ITrackRepository --output internal/adapters/repository/mocks \
    		--filename track.go --structname TrackRepository
	mockery --dir internal/ports --name ITokenProvider --output internal/adapters/auth/mocks \
		--filename auth.go --structname TokenProvider
	mockery --dir internal/ports --name IHashPasswordProvider --output internal/adapters/hash/mocks \
		--filename hash.go --structname HashPasswordProvider

test: 
	rm -rf allure-results
	go test -shuffle on \
		./internal/service/test \
		./internal/adapters/repository/postgres/test --parallel 8

allure:
	cp -R allure-reports/history allure-results
	rm -rf allure-reports
	allure generate allure-results -o allure-reports
	allure serve allure-results -p 4000

report: test allure

swagger:
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/web/main.go
	swagger2openapi docs/swagger.yaml -o docs/openapi3.yaml

