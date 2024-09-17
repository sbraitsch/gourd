package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"gourd/internal/views"
	"io"
	"net/http"
	"os"
)

var _ = Describe("Integration Tests", func() {

	Describe("Calling the base route", func() {
		It("should return index.html", func() {
			index, err := os.ReadFile("../internal/static/index.html")
			Expect(err).To(BeNil())

			resp, err := client.Get(fmt.Sprintf("http://localhost:%d/", testConfig.ServerPort))
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			body, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal(string(index)))
		})
	})

	Describe("Calling a protected route", func() {
		Context("without a token cookie", func() {
			It("should return the login page", func() {
				loginComponent := views.Login(testConfig.ApplicationTitle, testConfig.ApplicationSubtitle, testConfig.LogoPath)
				var buf bytes.Buffer
				loginComponent.Render(context.Background(), &buf)
				protectedPath := "api/content"

				resp, err := client.Get(fmt.Sprintf("http://localhost:%d/%s", testConfig.ServerPort, protectedPath))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(buf.String()))
			})
		})

		Context("with an invalid token cookie", func() {
			It("should return 404: Token not recognized", func() {
				protectedPath := "api/content"
				expectedBody := "Token not recognized\n"

				req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", testConfig.ServerPort, protectedPath), nil)
				Expect(err).To(BeNil())
				req.Header.Add("Cookie", "token=unrecognized")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(expectedBody))
			})
		})

		Context("with a valid token cookie", func() {
			It("should return 200 OK", func() {
				protectedPath := "api/content"

				req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", testConfig.ServerPort, protectedPath), nil)
				Expect(err).To(BeNil())
				req.Header.Add("Cookie", fmt.Sprintf("token=%s", testUserToken))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				log.Info().Msg(string(body))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Calling an admin route", func() {
		Context("without a token cookie", func() {
			It("should return the login page", func() {
				loginComponent := views.Login(testConfig.ApplicationTitle, testConfig.ApplicationSubtitle, testConfig.LogoPath)
				var buf bytes.Buffer
				loginComponent.Render(context.Background(), &buf)
				adminPath := "admin/generate"

				resp, err := client.Get(fmt.Sprintf("http://localhost:%d/%s", testConfig.ServerPort, adminPath))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(buf.String()))
			})
		})

		Context("with an invalid token cookie", func() {
			It("should return 404: Token not recognized", func() {
				adminPath := "admin/generate"
				expectedBody := "Token not recognized\n"

				req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", testConfig.ServerPort, adminPath), nil)
				Expect(err).To(BeNil())
				req.Header.Add("Cookie", "token=unrecognized")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(expectedBody))
			})
		})

		Context("with a non-admin token", func() {
			It("should return 403: Missing privileges", func() {
				adminPath := "admin/generate"
				expectedBody := "Missing privileges\n"

				req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", testConfig.ServerPort, adminPath), nil)
				Expect(err).To(BeNil())
				req.Header.Add("Cookie", fmt.Sprintf("token=%s", testUserToken))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(expectedBody))
			})
		})

		Context("with an admin token", func() {
			It("should return 200 OK", func() {
				adminPath := "admin/generate"

				req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", testConfig.ServerPort, adminPath), nil)
				Expect(err).To(BeNil())
				req.Header.Add("Cookie", fmt.Sprintf("token=%s", testAdminToken))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
