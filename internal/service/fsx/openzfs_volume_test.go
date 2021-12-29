package fsx_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/fsx"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tffsx "github.com/hashicorp/terraform-provider-aws/internal/service/fsx"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccFSxOpenzfsVolume_basic(t *testing.T) {
	var volume fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "fsx", regexp.MustCompile(`volume/fs-.+/fsvol-.+`)),
					resource.TestCheckResourceAttr(resourceName, "copy_tags_to_snapshots", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_compression_type", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.0.clients", "*"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.0.options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.0.options.0", "crossmnt"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_volume_id"),
					resource.TestCheckResourceAttr(resourceName, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_parentVolume(t *testing.T) {
	var volume, volume2 fsx.Volume
	var volumeId string
	resourceName := "aws_fsx_openzfs_volume.test"
	resourceName2 := "aws_fsx_openzfs_volume.test2"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeParentVolumeConfig(rName, rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume),
					testAccCheckFsxOpenzfsVolumeExists(resourceName2, &volume2),
					testAccCheckFsxOpenzfsVolumeGetId(resourceName, &volumeId),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "fsx", regexp.MustCompile(`volume/fs-.+/fsvol-.+`)),
					acctest.MatchResourceAttrRegionalARN(resourceName2, "arn", "fsx", regexp.MustCompile(`volume/fs-.+/fsvol-.+`)),
					resource.TestCheckResourceAttrPtr(resourceName2, "parent_volume_id", &volumeId),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_tags(t *testing.T) {
	var volume1, volume2, volume3 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeTags1Config(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeTags2Config(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccOpenzfsVolumeTags1Config(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume3),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume2, &volume3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_copyTags(t *testing.T) {
	var volume1, volume2 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeCopyTagsConfig(rName, "key1", "value1", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "copy_tags_to_snapshots", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeCopyTagsConfig(rName, "key1", "value1", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "copy_tags_to_snapshots", "false"),
				),
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_name(t *testing.T) {
	var volume1, volume2 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeBasicConfig(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_dataCompressionType(t *testing.T) {
	var volume1, volume2 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeDataCompressionConfig(rName, "ZSTD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "data_compression_type", "ZSTD"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeDataCompressionConfig(rName, "NONE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "data_compression_type", "NONE"),
				),
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_readOnly(t *testing.T) {
	var volume1, volume2 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeReadOnlyConfig(rName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "read_only", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeReadOnlyConfig(rName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "read_only", "true"),
				),
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_storageCapacity(t *testing.T) {
	var volume1, volume2 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeStorageCapacityConfig(rName, 30, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "storage_capacity_quota_gib", "30"),
					resource.TestCheckResourceAttr(resourceName, "storage_capacity_reservation_gib", "20"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeStorageCapacityConfig(rName, 40, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "storage_capacity_quota_gib", "40"),
					resource.TestCheckResourceAttr(resourceName, "storage_capacity_reservation_gib", "30"),
				),
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_nfsExports(t *testing.T) {
	var volume1, volume2 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeNFSExports1Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.0.clients", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.0.options.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.0.options.0", "async"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.0.options.1", "rw"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeNFSExports2Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nfs_exports.0.client_configurations.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "nfs_exports.0.client_configurations.*", map[string]string{
						"clients":   "10.0.1.0/24",
						"options.0": "async",
						"options.1": "rw",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "nfs_exports.0.client_configurations.*", map[string]string{
						"clients":   "*",
						"options.0": "sync",
						"options.1": "rw",
					}),
				),
			},
		},
	})
}

func TestAccFSxOpenzfsVolume_userAndGroupQuotas(t *testing.T) {
	var volume1, volume2 fsx.Volume
	resourceName := "aws_fsx_openzfs_volume.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxOpenzfsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenzfsVolumeUserAndGroupQuotas1Config(rName, 256),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume1),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.0.id", "10"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.0.storage_capacity_quota_gib", "256"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.0.type", "USER"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenzfsVolumeUserAndGroupQuotas2Config(rName, 128, 1024),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxOpenzfsVolumeExists(resourceName, &volume2),
					testAccCheckFsxOpenzfsVolumeNotRecreated(&volume1, &volume2),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.0.id", "10"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.0.storage_capacity_quota_gib", "128"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.0.type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.1.id", "20"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.1.storage_capacity_quota_gib", "1024"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.1.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.2.id", "5"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.2.storage_capacity_quota_gib", "1024"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.2.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.3.id", "100"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.3.storage_capacity_quota_gib", "128"),
					resource.TestCheckResourceAttr(resourceName, "user_and_group_quotas.3.type", "USER"),
				),
			},
		},
	})
}

func testAccCheckFsxOpenzfsVolumeExists(resourceName string, volume *fsx.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).FSxConn

		volume1, err := tffsx.FindVolumeByID(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		if volume == nil {
			return fmt.Errorf("FSx OpenZFS Volume (%s) not found", rs.Primary.ID)
		}

		*volume = *volume1

		return nil
	}
}

func testAccCheckFsxOpenzfsVolumeDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).FSxConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_fsx_openzfs_volume" {
			continue
		}

		volume, err := tffsx.FindVolumeByID(conn, rs.Primary.ID)
		if tfresource.NotFound(err) {
			continue
		}

		if volume != nil {
			return fmt.Errorf("FSx OpenZFS Volume (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckFsxOpenzfsVolumeGetId(resourceName string, volumeId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		*volumeId = rs.Primary.ID

		return nil
	}
}

func testAccCheckFsxOpenzfsVolumeNotRecreated(i, j *fsx.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.VolumeId) != aws.StringValue(j.VolumeId) {
			return fmt.Errorf("FSx OpenZFS Volume (%s) recreated", aws.StringValue(i.VolumeId))
		}

		return nil
	}
}

func testAccCheckFsxOpenzfsVolumeRecreated(i, j *fsx.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.VolumeId) == aws.StringValue(j.VolumeId) {
			return fmt.Errorf("FSx OpenZFS Volume (%s) not recreated", aws.StringValue(i.VolumeId))
		}

		return nil
	}
}

func testAccOpenzfsVolumeBaseConfig(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test1" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_fsx_openzfs_file_system" "test" {
  storage_capacity    = 64
  subnet_ids          = [aws_subnet.test1.id]
  deployment_type     = "SINGLE_AZ_1"
  throughput_capacity = 64
  }

`, rName))
}

func testAccOpenzfsVolumeBasicConfig(rName string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
}
`, rName))
}

func testAccOpenzfsVolumeParentVolumeConfig(rName, rName2 string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
}

resource "aws_fsx_openzfs_volume" "test2" {
  name             = %[2]q
  parent_volume_id = aws_fsx_openzfs_volume.test.id
  }
`, rName, rName2))
}

func testAccOpenzfsVolumeTags1Config(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccOpenzfsVolumeTags2Config(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
	

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccOpenzfsVolumeCopyTagsConfig(rName, tagKey1, tagValue1, copyTags string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name                   = %[1]q
  parent_volume_id       = aws_fsx_openzfs_file_system.test.root_volume_id
  copy_tags_to_snapshots = %[4]s

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1, copyTags))
}

func testAccOpenzfsVolumeDataCompressionConfig(rName, dType string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name                  = %[1]q
  parent_volume_id      = aws_fsx_openzfs_file_system.test.root_volume_id
  data_compression_type = %[2]q
}
`, rName, dType))
}

func testAccOpenzfsVolumeReadOnlyConfig(rName, readOnly string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
  read_only        = %[2]s
}
`, rName, readOnly))
}

func testAccOpenzfsVolumeStorageCapacityConfig(rName string, storageQuota, storageReservation int) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name                             = %[1]q
  parent_volume_id                 = aws_fsx_openzfs_file_system.test.root_volume_id
  storage_capacity_quota_gib       = %[2]d
  storage_capacity_reservation_gib = %[3]d
}
`, rName, storageQuota, storageReservation))
}

func testAccOpenzfsVolumeNFSExports1Config(rName string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
  nfs_exports {
    client_configurations {
      clients = "10.0.1.0/24"
      options = ["async", "rw"]
    }
  }
    
}
`, rName))
}

func testAccOpenzfsVolumeNFSExports2Config(rName string) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
    nfs_exports {
      client_configurations {
        clients = "10.0.1.0/24"
        options = ["async", "rw"]
      }
      client_configurations {
        clients = "*"
        options = ["sync", "rw"]
      }
    }
}
`, rName))
}

func testAccOpenzfsVolumeUserAndGroupQuotas1Config(rName string, quotaSize int) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
  user_and_group_quotas {
	id                         = 10
	storage_capacity_quota_gib = %[2]d
	type                       = "USER"
  }
}
`, rName, quotaSize))
}

func testAccOpenzfsVolumeUserAndGroupQuotas2Config(rName string, userQuota, groupQuota int) string {
	return acctest.ConfigCompose(testAccOpenzfsVolumeBaseConfig(rName), fmt.Sprintf(`
resource "aws_fsx_openzfs_volume" "test" {
  name             = %[1]q
  parent_volume_id = aws_fsx_openzfs_file_system.test.root_volume_id
  user_and_group_quotas {
    id                         = 10
    storage_capacity_quota_gib = %[2]d
    type                       = "USER"
  }
  user_and_group_quotas {
    id                         = 20
    storage_capacity_quota_gib = %[3]d
    type                       = "GROUP"
  }
  user_and_group_quotas {
    id                         = 5
    storage_capacity_quota_gib = %[3]d
    type                       = "GROUP"
  }
  user_and_group_quotas {
    id                         = 100
    storage_capacity_quota_gib = %[2]d
    type                       = "USER"
  }
}
`, rName, userQuota, groupQuota))
}
