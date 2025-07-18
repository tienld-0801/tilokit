package scaffold

import (
	"fmt"
	"os"
)

func GenerateLaravel(projectName string) error {
	fmt.Println("🚧 Đang tạo project React:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("không thể tạo thư mục: %w", err)
	}

	return nil
}
