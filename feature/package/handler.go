package package_bill

import (
	"kubik-rental/config"
	"kubik-rental/dto"
	"kubik-rental/entity"
	"kubik-rental/pkg"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreatePackage(c *fiber.Ctx) error {
	var input CreatePackageRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{Message: err.Error(), Status: fiber.StatusBadRequest})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		errMap := pkg.FormatValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Validation failed", "Status": fiber.StatusBadRequest, "error": errMap})
	}

	newPkg := entity.Package{
		Name:     input.Name,
		Duration: input.Duration,
		Price:    input.Price,
	}

	if err := config.DB.Create(&newPkg).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{Message: err.Error(), Status: fiber.StatusBadRequest})
	}

	res := PackageResponse{
		ID:        newPkg.ID,
		Name:      newPkg.Name,
		Duration:  newPkg.Duration,
		Price:     newPkg.Price,
		CreatedAt: newPkg.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response[PackageResponse]{Message: "Package created successfully", Status: fiber.StatusCreated, Data: res})
}

func GetPackageByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var pkgData entity.Package
	if err := config.DB.First(&pkgData, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{Message: "Package not found", Status: fiber.StatusNotFound})
	}

	res := PackageResponse{
		ID:        pkgData.ID,
		Name:      pkgData.Name,
		Duration:  pkgData.Duration,
		Price:     pkgData.Price,
		CreatedAt: pkgData.CreatedAt,
		UpdatedAt: pkgData.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[PackageResponse]{Message: "Package found", Status: fiber.StatusOK, Data: res})
}

func GetAllPackages(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	keyword := c.Query("keyword", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if err := config.DB.Model(&entity.Package{}).Where("name LIKE ?", "%"+keyword+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{Message: "Failed to count packages", Status: fiber.StatusBadRequest})
	}

	var packages []entity.Package
	if err := config.DB.Limit(limit).Offset(offset).Where("name LIKE ?", "%"+keyword+"%").Find(&packages).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{Message: err.Error(), Status: fiber.StatusBadRequest})
	}

	result := make([]PackageResponse, 0)

	for _, p := range packages {
		result = append(result, PackageResponse{
			ID:        p.ID,
			Name:      p.Name,
			Duration:  p.Duration,
			Price:     p.Price,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	meta := dto.Meta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[[]PackageResponse]{Message: "Packages retrieved successfully", Status: fiber.StatusOK, Data: result, Meta: &meta})
}

func UpdatePackage(c *fiber.Ctx) error {
	id := c.Params("id")
	var pkgData entity.Package
	if err := config.DB.First(&pkgData, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Response[any]{Message: "Package not found", Status: fiber.StatusNotFound})
	}

	var input UpdatePackageRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{Message: err.Error(), Status: fiber.StatusBadRequest})
	}

	if err := pkg.Validate.Struct(input); err != nil {
		errMap := pkg.FormatValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Validation failed", "Status": fiber.StatusBadRequest, "error": errMap})
	}

	pkgData.Name = input.Name
	pkgData.Duration = input.Duration
	pkgData.Price = input.Price

	if err := config.DB.Save(&pkgData).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response[any]{Message: err.Error(), Status: fiber.StatusBadRequest})
	}

	res := PackageResponse{
		ID:        pkgData.ID,
		Name:      pkgData.Name,
		Duration:  pkgData.Duration,
		Price:     pkgData.Price,
		CreatedAt: pkgData.CreatedAt,
		UpdatedAt: pkgData.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[PackageResponse]{Message: "Package updated successfully", Status: fiber.StatusOK, Data: res})
}

func DeletePackage(c *fiber.Ctx) error {
	id := c.Params("id")
	result := config.DB.Delete(&entity.Package{}, id)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Package not found"})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response[any]{Message: "Package deleted successfully", Status: fiber.StatusOK})
}
