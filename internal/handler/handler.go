package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"nat.vg/narwhal/internal/model"
	"nat.vg/narwhal/internal/service"
)

type Handler struct {
	svc service.Service
}

func NarwhalHandler(e *echo.Echo, ur service.Service) {
	h := &Handler{
		svc: ur,
	}
	e.GET("/node/resource_usage", h.GetNodeResourceUsage)
	e.GET("/containers", h.ListContainers)
	e.GET("/containers/:id/stats", h.GetContainerStats)
	e.POST("/containers/create_simple", h.CreateContainerSimple)
	e.POST("/containers/create", h.CreateContainer)
	e.GET("/containers/:id", h.GetContainer)
	e.GET("/containers/:id/stop", h.StopContainer)
	e.GET("/containers/:id/start", h.StartContainer)
	e.DELETE("/containers/:id", h.RemoveContainer)

	e.GET("/containers/:id/logs", h.GetContainerLogs)

}

func (h *Handler) GetNodeResourceUsage(c echo.Context) error {
	// get the resource usage
	v, _ := mem.VirtualMemory()
	cCounts, _ := cpu.Counts(true)
	c, _ :=  cpu.Percent(time.Duration(1)*time.Second), true)

	// return an error
	return echo.NewHTTPError(http.StatusNotImplemented, fmt.Errorf("not implemented"))
}

func (h *Handler) ListContainers(c echo.Context) error {
	containers, err := h.svc.List(c.Request().Context())
	if err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, containers)
}

func (h *Handler) GetContainer(c echo.Context) error {
	id := c.Param("id")
	container, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, container)
}

func (h *Handler) GetContainerStats(c echo.Context) error {
	id := c.Param("id")
	container_info, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var basicResUsage *model.BasicResourceUsage

	if container_info.State.Running {
		res_usage, err := h.svc.GetStats(c.Request().Context(), id)
		if err != nil {
			// log error
			fmt.Fprintln(os.Stderr, err.Error())
			// return error
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		basicResUsage = &model.BasicResourceUsage{
			RamUsage: float64(res_usage.MemoryStats.Usage),
			CpuUsage: float64(res_usage.CPUStats.CPUUsage.TotalUsage)}
	}

	stats := model.BasicContainerStatistics{
		Name:  container_info.Name,
		Id:    container_info.ID,
		Image: container_info.Config.Image,
		//DiskUsage: res_usage.BlkioStats.
		ResourceUsage: basicResUsage,
		RamTotal:      float64(container_info.HostConfig.Memory),
		CpuTotal:      float64(container_info.HostConfig.NanoCPUs),
		StartCmd:      strings.Join(container_info.Config.Cmd, " "),
		IsRunning:     container_info.State.Running,
	}
	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) CreateContainerSimple(c echo.Context) error {
	req := new(model.CreateSimpleContainerRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	container, err := h.svc.CreateSimple(c.Request().Context(), *req)
	if err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, container)
}

func (h *Handler) CreateContainer(c echo.Context) error {
	req := new(model.CreateContainerRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := h.svc.Create(c.Request().Context(), *req); err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, req)
}

func (h *Handler) StopContainer(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.Stop(c.Request().Context(), id); err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, id)
}

func (h *Handler) StartContainer(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.Start(c.Request().Context(), id); err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, id)
}

func (h *Handler) RemoveContainer(c echo.Context) error {
	id := c.Param("id")
	force := c.QueryParam("force") == "true"
	removeVolumes := c.QueryParam("remove_volumes") == "true"
	if err := h.svc.Remove(c.Request().Context(), id, force, removeVolumes); err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, id)
}

func (h *Handler) GetContainerLogs(c echo.Context) error {
	id := c.Param("id")
	// to int
	lines, err := strconv.Atoi(c.QueryParam("lines"))
	if err != nil {
		lines = 50
	}
	since := c.QueryParam("since")
	stream := c.QueryParam("stream") == "true"

	logs, err := h.svc.GetLogs(c.Request().Context(), id, lines, since, stream)
	if err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	defer logs.Close()
	_, err = io.Copy(c.Response().Writer, logs)
	if err != nil {
		// log error
		fmt.Fprintln(os.Stderr, err.Error())
		// return error
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.Request().Body.Close()
}
