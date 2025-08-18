package handler

import (
	"net/http"
	"strconv"

	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/api/response"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/service"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service service.PatientService
}

func NewPatientHandler(s service.PatientService) *PatientHandler {
	return &PatientHandler{service: s}
}

func (h *PatientHandler) GetPatients(c *gin.Context) {
	patients, err := h.service.GetPatients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: "failed to fetch patients",
			Details: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Patients fetched successfully",
		Data:    patients,
	})
}

func (h *PatientHandler) GetPatientByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "Invalid patient ID",
			Details: err.Error(),
		})
		return
	}

	patient, err := h.service.GetPatientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "Patient not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Patient retrieved successfully",
		Data:    patient,
	})
}

func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var patient domain.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	if _, err := h.service.CreatePatient(patient); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: "Failed to create patient",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse{
		Message: "Patient created successfully",
		Data:    patient,
	})
}

func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "Invalid patient ID",
			Details: err.Error(),
		})
		return
	}

	var patient domain.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	patient.ID = uint(id)
	if err := h.service.UpdatePatient(&patient); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: "Failed to update patient",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Patient updated successfully",
		Data:    patient,
	})
}

func (h *PatientHandler) DeletePatient(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "Invalid patient ID",
			Details: err.Error(),
		})
		return
	}

	if err := h.service.DeletePatient(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: "Failed to delete patient",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Patient deleted successfully",
	})
}
