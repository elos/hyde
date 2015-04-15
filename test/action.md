Action
------

An action is something a user can be doing. An action should not be confused with an [event](event.md), which is anything that can happen to a user. For example, you can imagine synchronous events - a user can be studying, and at the library. The events could be split: studying would mirror being at the library. The action would only ever be studying. Perhaps "walking to the library" would be another.

An action is owned by a [user](user.md), has a name, a start time and an end time, a boolean field marking it's completion, and an associated actionable, which "created" the task.


### Implementation
| Kind   | Concerns        |
| ------ | --------------- |
| action | [user](user.md) |

### Transfer Structure
| Attribute      | JSON Key      | Type                | Access    |
| -------------- | ------------- | ------------------- | --------- |
| ID             | id              | bson.ObjectId     | Public    |
| CreatedAt      | created_at      | RFC 3339 DateTime | Personal  |
| UpdatedAt      | updated_at      | RFC 3339 DateTime | Personal  |
| Name           | name            | String            | Personal  |
| StartTime      | start_time      | RFC 3339 DateTime | Personal  |
| EndTime        | end_time        | RFC 3339 DateTime | Personal  |
| UserID         | user_id         | bson.ObjectId     | Personal  |
| Completed      | completed       | boolean           | Personal  |
| ActionableKind | actionable_kind | string            | Personal  |
| ActionableID   | actionable_id   | bson.ObjectId     | Personal  |
