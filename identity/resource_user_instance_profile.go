package identity

import (
	"context"
	"fmt"

	"github.com/databrickslabs/databricks-terraform/common"
	"github.com/databrickslabs/databricks-terraform/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceUserInstanceProfile binds user and instance profile
func ResourceUserInstanceProfile() *schema.Resource {
	return util.NewPairID("user_id", "instance_profile_id").BindResource(util.BindResource{
		CreateContext: func(ctx context.Context, userID, roleARN string, c *common.DatabricksClient) error {
			err := validateInstanceProfileARN(roleARN)
			if err != nil {
				return err
			}
			return NewUsersAPI(c).Patch(userID, scimPatchRequest("add", "roles", roleARN))
		},
		ReadContext: func(ctx context.Context, userID, roleARN string, c *common.DatabricksClient) error {
			user, err := NewUsersAPI(c).read(userID)
			if err == nil && !user.HasRole(roleARN) {
				return common.NotFound("User has no role")
			}
			return err
		},
		DeleteContext: func(ctx context.Context, userID, roleARN string, c *common.DatabricksClient) error {
			return NewUsersAPI(c).Patch(userID, scimPatchRequest(
				"remove", fmt.Sprintf(`roles[value eq "%s"]`, roleARN), ""))
		},
	})
}
