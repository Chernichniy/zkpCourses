package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"net/http"

	"github.com/Chernichniy/zkpCourses/goScripts/Lagrangia"
	qap "github.com/Chernichniy/zkpCourses/goScripts/QAP"
	r1cs "github.com/Chernichniy/zkpCourses/goScripts/R1CS"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

type lagrangiaOutput struct {
	Polynomials                   map[string]string `json:"Polynomials"`
	FullPolinomialMultiPointBasis string            `json:"PolynomialsBaric"`
	Basises                       map[int]string    `json:"Basises"`
}

type qapOutput struct {
	Module          int        `json:"Module"`
	R1CS            r1csOutput `json:"R1CSParam"`
	QAPMatrixA      [][]int    `json:"QAPMatrixA"`
	QAPMatrixB      [][]int    `json:"QAPMatrixB"`
	QAPMatrixC      [][]int    `json:"QAPMatrixC"`
	FinallVector    []int      `json:"FinallVector"`
	VanishVector    []int      `json:"VanishVector"`
	QuiotientVector []int      `json:"Quiotient"`
}

type r1csOutput struct {
	WitnessFormal       []string `json:"WitnessFormalForm"`
	WitnesFormalChecker []string `json:"WitnessFormalChecker"`
	WitnessNumber       []int    `json:"WitnessNumberForm"`
	WitnesNumberChecker []string `json:"WitnessNumberChecker"`
	MatrixA             [][]int  `json:"MatrixA"`
	MatrixB             [][]int  `json:"MatrixB"`
	MatrixC             [][]int  `json:"MatrixC"`
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	s.logger.Info("Starting API server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {

	s.router.HandleFunc("/r1cs/{function}/{roots}", s.handleR1CS()).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/Lagrangia/{roots}", s.handleLagrangia()).Methods("GET")
	s.router.HandleFunc("/QAP/{function}/{roots}/{mod}", s.handleQAP())

}

func (s *APIServer) handleR1CS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		inputFuncAndRoots := mux.Vars(r)
		fmt.Println(inputFuncAndRoots)

		if inputFuncAndRoots["roots"] != "1" { // Value "1" for 1 time page download

			r1cs.Start(inputFuncAndRoots["function"], inputFuncAndRoots["roots"])
		} else {
			r1cs.Start("x ^ 3 + x + 5", "x = 2 y = 15")
		}

		newR1CSOutput := &r1csOutput{
			WitnessFormal:       r1cs.ReturnWitnessFormal(),
			WitnesFormalChecker: r1cs.ReturnWitnessFormalChecker(),
			WitnessNumber:       r1cs.ReturnWitnessNumbers(),
			WitnesNumberChecker: r1cs.ReturnWitnessNumberChecker(),
			MatrixA:             r1cs.ReturnVectorsA(),
			MatrixB:             r1cs.ReturnVectorsB(),
			MatrixC:             r1cs.ReturnVectorsC(),
		}

		var returnJsonR1CS bytes.Buffer

		encodedJsonStruct := json.NewEncoder(&returnJsonR1CS)

		encodedJsonStruct.SetIndent("", "")
		encodedJsonStruct.Encode(&newR1CSOutput)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		io.WriteString(w, returnJsonR1CS.String())

		r1cs.ClearAllVar()

	}
}

// Здесь необходимо ретурны сделать так, что бы приходили разные данные в зависимости от количек точек интерполяции
func (s *APIServer) handleLagrangia() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		roots := mux.Vars(r)
		fmt.Println(roots)

		if roots["roots"] != "1" { // Value "1" for 1 time page download
			Lagrangia.Start(roots["roots"])
		} else {
			Lagrangia.Start("x = 1 x1 = 2 x2 = 3 y = 3 y1 = 2 y2 = 4")
		}

		newLagrangiaOutput := &lagrangiaOutput{
			Polynomials:                   Lagrangia.ReturnMultiPolynomial(),
			FullPolinomialMultiPointBasis: Lagrangia.ReturnMultiPolynomialBarric(),
			Basises:                       Lagrangia.ReturnBaricBasis(),
		}

		var returnJsonLagrangia bytes.Buffer

		encodedJsonStruct := json.NewEncoder(&returnJsonLagrangia)

		encodedJsonStruct.SetIndent("", "\t")
		encodedJsonStruct.Encode(&newLagrangiaOutput)

		io.WriteString(w, returnJsonLagrangia.String())

		Lagrangia.ClearAllVar()
	}
}

func (s *APIServer) handleQAP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		inputFuncAndRoots := mux.Vars(r)
		fmt.Println(inputFuncAndRoots)

		moduleInput, _ := strconv.Atoi(inputFuncAndRoots["mod"])

		fmt.Println(inputFuncAndRoots["function"], inputFuncAndRoots["roots"], inputFuncAndRoots["mod"])

		if inputFuncAndRoots["roots"] != "1" { // Value "1" for 1 time page download
			qap.Start(inputFuncAndRoots["function"], inputFuncAndRoots["roots"], moduleInput)
		} else {
			qap.Start("x ^ 3 + x + 5", "x = 2 y = 15", 11)
		}

		newQAPOutput := &qapOutput{
			Module: moduleInput,
			R1CS: r1csOutput{
				r1cs.ReturnWitnessFormal(),
				r1cs.ReturnWitnessFormalChecker(),
				r1cs.ReturnWitnessNumbers(),
				r1cs.ReturnWitnessNumberChecker(),
				r1cs.ReturnVectorsA(),
				r1cs.ReturnVectorsA(),
				r1cs.ReturnVectorsA()},
			QAPMatrixA:      qap.QAPVectAReturn(),
			QAPMatrixB:      qap.QAPVectBReturn(),
			QAPMatrixC:      qap.QAPVectCReturn(),
			FinallVector:    qap.QAPVectFinallReturn(),
			VanishVector:    qap.QAPVanishVectorReturn(),
			QuiotientVector: qap.QAPQuiotientOfFinallFructReturn(),
		}

		var returnJsonQAP bytes.Buffer

		encodedJsonStruct := json.NewEncoder(&returnJsonQAP)

		encodedJsonStruct.SetIndent("", "")
		encodedJsonStruct.Encode(&newQAPOutput)

		io.WriteString(w, returnJsonQAP.String())

		qap.ClearAllVar()

	}
}
