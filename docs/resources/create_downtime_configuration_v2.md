---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "synthetics_create_downtime_configuration_v2 Resource - synthetics"
subcategory: ""
description: |-
  
---

# synthetics_create_downtime_configuration_v2 (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `downtime_configuration` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--downtime_configuration))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--downtime_configuration"></a>
### Nested Schema for `downtime_configuration`

Required:

- `end_time` (String)
- `name` (String)
- `rule` (String)
- `start_time` (String)
- `test_ids` (List of Number)

Optional:

- `description` (String)

Read-Only:

- `created_at` (String)
- `id` (Number) The ID of this resource.
- `status` (String)
- `tests_updated_at` (String)
- `updated_at` (String)


