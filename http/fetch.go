package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

// fetchRemoteResource: 범용적인 데이터 가져오기 함수
func fetchRemoteResource(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %s", r.Status)
	}

	return io.ReadAll(r.Body)
}

// startTestPackageServer: 라우팅이 포함된 개선된 테스트 서버
func startTestPackageServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/packages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"name":"package1","version":"1.0.0"},{"name":"package2","version":"2.0.0"}]`))
		case http.MethodPost:
			packageRegHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return httptest.NewServer(mux)
}

// pkgData & pkgRegisterResult
type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type pkgRegisterResult struct {
	Id string `json:"id"`
}

// fetchPackageData: JSON 역직렬화 처리
func fetchPackageData(url string) ([]pkgData, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("unexpected content type")
	}

	var packages []pkgData
	if err := json.NewDecoder(r.Body).Decode(&packages); err != nil {
		return nil, err
	}
	return packages, nil
}

// downloadToFile: 메모리 효율적인 스트림 방식 저장
func downloadToFile(url string, filePath string) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", r.Status)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	return err
}

// registerPackageData: POST 요청 및 응답 처리
func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return pkgRegisterResult{}, err
	}

	r, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return pkgRegisterResult{}, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r.Body)
		return pkgRegisterResult{}, fmt.Errorf("server error: %s", string(body))
	}

	var result pkgRegisterResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return pkgRegisterResult{}, err
	}
	return result, nil
}

// packageRegHandler: 핸들러 내부 switch문 적용
func packageRegHandler(w http.ResponseWriter, r *http.Request) {
	var p pkgData
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.Name == "" || p.Version == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	d := pkgRegisterResult{Id: fmt.Sprintf("%s-%s", p.Name, p.Version)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}
