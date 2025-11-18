package api

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/util/serializer"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// agent 图片上传
func agentLogoUpload(w http.ResponseWriter, r *http.Request) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmServApi)
	log.Info("do")
	if r.Method != http.MethodPost {
		serializer.StringErrorRes(w, errors.New("Method Not Allowed"))
		return
	}

	// 解析 multipart 表单（限制内存使用 1MB，超过部分存临时文件）
	err := r.ParseMultipartForm(1 << 20) // 1MB
	if err != nil {
		log.WithError(err).Error("ParseMultipartForm")
		if err.Error() == "request too large" {
			serializer.StringErrorRes(w, errors.New("上传文件大小超过限制(1MB)"))
		} else {
			serializer.StringErrorRes(w, errors.New("解析form出错"))
		}
		return
	}

	// 获取上传的文件（表单字段名为 "image"）
	file, handler, err := r.FormFile("image")
	if err != nil {
		log.WithError(err).Error("FormFile")
		serializer.StringErrorRes(w, errors.New("获取图片文件出错"))
		return
	}
	defer file.Close()

	// 检查文件类型（可选：限制为图片格式）
	ext := filepath.Ext(handler.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		serializer.StringErrorRes(w, errors.New("Only JPG/JPEG/PNG/GIF images are allowed"))
		return
	}

	// 确保上传目录存在（./uploads）
	uploadDir := "./uploads"
	if err = os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.WithError(err).Error("os.MkdirAll")
		serializer.StringErrorRes(w, errors.New("创建目录失败"))
		return
	}

	// 生成保存路径（直接使用原文件名，生产环境建议使用 UUID 或时间戳重命名）
	fileExt := filepath.Ext(handler.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	savePath := filepath.Join(uploadDir, newFileName)
	fileWebPath := fmt.Sprintf("%s/logos/%s", core.BaseConf.LocalServerUrl, newFileName)
	dst, err := os.Create(savePath)
	if err != nil {
		log.WithError(err).Error("os.Create")
		serializer.StringErrorRes(w, errors.New("创建文件失败"))
		return
	}
	defer dst.Close()

	// 复制文件内容到目标路径
	size, err := io.Copy(dst, file)
	if err != nil {
		log.WithError(err).Error("io.Copy")
		serializer.StringErrorRes(w, errors.New("复制文件失败"))
		return
	}

	// 返回成功响应
	response := UploadResponse{
		Filename: fileWebPath,
		Size:     size,
	}

	serializer.SuccessRes(w, response)
}

type UploadResponse struct {
	Filename string `json:"filename"` // 保存的文件名
	Size     int64  `json:"size"`     // 文件大小（字节）
}
