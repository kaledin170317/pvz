package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Ошибка выполнения %s: %v\n", name, err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Генерация моков...")

	// переход в корень проекта
	root, _ := filepath.Abs("../..")
	if err := os.Chdir(root); err != nil {
		fmt.Println("Не удалось перейти в корень:", err)
		os.Exit(1)
	}

	run("mockgen", "-source=internal/domain/usecases/pvz_usecase.go", "-destination=internal/domain/usecases/mocks/pvz_usecase_mock.go", "-package=mocks")
	run("mockgen", "-source=internal/domain/usecases/reception_usecase.go", "-destination=internal/domain/usecases/mocks/reception_usecase_mock.go", "-package=mocks")
	run("mockgen", "-source=internal/domain/usecases/product_usecase.go", "-destination=internal/domain/usecases/mocks/product_usecase_mock.go", "-package=mocks")
	run("mockgen", "-source=internal/domain/usecases/user_usecases.go", "-destination=internal/domain/usecases/mocks/user_usecase_mock.go", "-package=mocks")
	run("mockgen", "-source=internal/adapters/db/product_repo.go", "-destination=internal/adapters/db/mocks/product_repo_mock.go", "-package=mocks")
	run("mockgen", "-source=internal/adapters/db/pvz_repo.go", "-destination=internal/adapters/db/mocks/pvz_repo_mock.go", "-package=mocks")
	run("mockgen", "-source=internal/adapters/db/reception_repo.go", "-destination=internal/adapters/db/mocks/reception_repo_mock.go", "-package=mocks")
	run("mockgen", "-source=internal/adapters/db/user_repo.go", "-destination=internal/adapters/db/mocks/user_repo_mock.go", "-package=mocks")

	fmt.Println("Моки сгенерированы")
}
