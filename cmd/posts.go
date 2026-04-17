package cmd

import (
	"fmt"
	"strings"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	postID                       string
	postListLimit                int
	postListCursor               string
	postListIDs                  string
	postListStatus               string
	postListChannel              string
	postListLinkedInProfileID    string
	postListXProfileID           string
	postListInstagramProfileID   string
	postListLabelIDs             string
	postListPublishDateMin       string
	postListPublishDateMax       string
	postListCreatedAtMin         string
	postListCreatedAtMax         string
	postListSortBy               string
	postListSortOrder            string
	postListAll                  bool
	postCreateTitle              string
	postCreatePublishAt          string
	postCreateStatus             string
	postCreateLabelIDs           string
	postCreateCampaignID         string
	postCreateNotes              string
	postCreateBodyJSON           string
	postCreateBodyFile           string
	postUpdateTitle              string
	postUpdatePublishAt          string
	postUpdateStatus             string
	postUpdateLabelIDs           string
	postUpdateCampaignID         string
	postUpdateNotes              string
	postUpdateBodyJSON           string
	postUpdateBodyFile           string
	postSchedulePublishAt        string
)

func init() {
	rootCmd.AddCommand(postCmd)
	postCmd.AddCommand(postListCmd)
	postCmd.AddCommand(postGetCmd)
	postCmd.AddCommand(postCreateCmd)
	postCmd.AddCommand(postUpdateCmd)
	postCmd.AddCommand(postArchiveCmd)
	postCmd.AddCommand(postUnarchiveCmd)
	postCmd.AddCommand(postScheduleCmd)
	postCmd.AddCommand(postUnscheduleCmd)

	postListCmd.Flags().IntVar(&postListLimit, "limit", 0, "Max posts to return (1-100)")
	postListCmd.Flags().StringVar(&postListCursor, "cursor", "", "Pagination cursor from a prior response")
	postListCmd.Flags().StringVar(&postListIDs, "ids", "", "Fetch specific posts by UUIDs (comma-separated)")
	postListCmd.Flags().StringVar(&postListStatus, "status", "", "Filter by post status")
	postListCmd.Flags().StringVar(&postListChannel, "channel", "", "Filter by channel")
	postListCmd.Flags().StringVar(&postListLinkedInProfileID, "linkedin-profile-id", "", "Filter by LinkedIn profile ID")
	postListCmd.Flags().StringVar(&postListXProfileID, "x-profile-id", "", "Filter by X/Twitter profile ID")
	postListCmd.Flags().StringVar(&postListInstagramProfileID, "instagram-profile-id", "", "Filter by Instagram profile ID")
	postListCmd.Flags().StringVar(&postListLabelIDs, "label-ids", "", "Filter by label IDs (comma-separated)")
	postListCmd.Flags().StringVar(&postListPublishDateMin, "publish-date-min", "", "Filter posts scheduled on or after this date")
	postListCmd.Flags().StringVar(&postListPublishDateMax, "publish-date-max", "", "Filter posts scheduled on or before this date")
	postListCmd.Flags().StringVar(&postListCreatedAtMin, "created-at-min", "", "Filter posts created on or after this date")
	postListCmd.Flags().StringVar(&postListCreatedAtMax, "created-at-max", "", "Filter posts created on or before this date")
	postListCmd.Flags().StringVar(&postListSortBy, "sort-by", "", "Field to sort by")
	postListCmd.Flags().StringVar(&postListSortOrder, "sort-order", "", "Sort order (asc or desc)")
	postListCmd.Flags().BoolVar(&postListAll, "all", false, "Auto-paginate and return all results")

	postGetCmd.Flags().StringVar(&postID, "id", "", "Post ID (UUID)")
	postGetCmd.MarkFlagRequired("id")

	postCreateCmd.Flags().StringVar(&postCreateTitle, "title", "", "Post title")
	postCreateCmd.Flags().StringVar(&postCreatePublishAt, "publish-at", "", "Publish datetime UTC (ISO 8601)")
	postCreateCmd.Flags().StringVar(&postCreateStatus, "status", "", "Post status (Tentative, ToDo, InProgress, ForReview, Blocked, Finalized, Scheduled)")
	postCreateCmd.Flags().StringVar(&postCreateLabelIDs, "label-ids", "", "Comma-separated label IDs")
	postCreateCmd.Flags().StringVar(&postCreateCampaignID, "campaign-id", "", "Campaign ID")
	postCreateCmd.Flags().StringVar(&postCreateNotes, "notes", "", "Post notes")
	postCreateCmd.Flags().StringVar(&postCreateBodyJSON, "body-json", "", "Full JSON body (individual flags override matching keys; use for nested channel configs)")
	postCreateCmd.Flags().StringVar(&postCreateBodyFile, "body-file", "", "Path to JSON body file (or - for stdin)")

	postUpdateCmd.Flags().StringVar(&postID, "id", "", "Post ID (UUID)")
	postUpdateCmd.Flags().StringVar(&postUpdateTitle, "title", "", "Post title")
	postUpdateCmd.Flags().StringVar(&postUpdatePublishAt, "publish-at", "", "Publish datetime UTC (ISO 8601)")
	postUpdateCmd.Flags().StringVar(&postUpdateStatus, "status", "", "Post status")
	postUpdateCmd.Flags().StringVar(&postUpdateLabelIDs, "label-ids", "", "Comma-separated label IDs")
	postUpdateCmd.Flags().StringVar(&postUpdateCampaignID, "campaign-id", "", "Campaign ID")
	postUpdateCmd.Flags().StringVar(&postUpdateNotes, "notes", "", "Post notes")
	postUpdateCmd.Flags().StringVar(&postUpdateBodyJSON, "body-json", "", "Full JSON body (individual flags override matching keys when set)")
	postUpdateCmd.Flags().StringVar(&postUpdateBodyFile, "body-file", "", "Path to JSON body file (or - for stdin)")
	postUpdateCmd.MarkFlagRequired("id")

	postArchiveCmd.Flags().StringVar(&postID, "id", "", "Post ID (UUID)")
	postArchiveCmd.MarkFlagRequired("id")
	postUnarchiveCmd.Flags().StringVar(&postID, "id", "", "Post ID (UUID)")
	postUnarchiveCmd.MarkFlagRequired("id")

	postScheduleCmd.Flags().StringVar(&postID, "id", "", "Post ID (UUID)")
	postScheduleCmd.Flags().StringVar(&postSchedulePublishAt, "publish-at", "", "Publish datetime UTC (ISO 8601); omit to reschedule using existing publishAt")
	postScheduleCmd.MarkFlagRequired("id")

	postUnscheduleCmd.Flags().StringVar(&postID, "id", "", "Post ID (UUID)")
	postUnscheduleCmd.MarkFlagRequired("id")
}

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Manage posts",
	Long:  "Create, list, update, schedule, archive, and retrieve posts across LinkedIn, X, and Instagram channels.",
}

var postListCmd = &cobra.Command{
	Use:   "list",
	Short: "List posts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Help text advertises 1-100; enforce locally so the flag's
		// contract, runtime behavior, and API constraints agree without
		// making a round-trip for an obviously-invalid value.
		if postListLimit < 0 || postListLimit > 100 {
			return fmt.Errorf("--limit must be between 1 and 100")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		params := models.ListPostsParams{
			Limit:              postListLimit,
			Cursor:             postListCursor,
			IDs:                postListIDs,
			Status:             postListStatus,
			Channel:            postListChannel,
			LinkedInProfileID:  postListLinkedInProfileID,
			XProfileID:         postListXProfileID,
			InstagramProfileID: postListInstagramProfileID,
			LabelIDs:           postListLabelIDs,
			PublishDateMin:     postListPublishDateMin,
			PublishDateMax:     postListPublishDateMax,
			CreatedAtMin:       postListCreatedAtMin,
			CreatedAtMax:       postListCreatedAtMax,
			SortBy:             postListSortBy,
			SortOrder:          postListSortOrder,
		}
		svc := api.NewPostService(c)
		if postListAll {
			posts, err := svc.ListAll(params)
			if err != nil {
				return err
			}
			return printResult(posts)
		}
		resp, err := svc.List(params)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var postGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a post by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		p, err := api.NewPostService(c).Get(postID)
		if err != nil {
			return err
		}
		return printResult(p)
	},
}

var postCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a post",
	Long:  "Create a post. Use --body-json or --body-file to pass the full request including nested linkedIn/x/instagram channel configs. Individual flags override matching top-level keys in the body when provided.",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := parseBodyJSON(postCreateBodyJSON, postCreateBodyFile)
		if err != nil {
			return err
		}
		if body == nil {
			body = map[string]interface{}{}
		}
		if postCreateTitle != "" {
			body["title"] = postCreateTitle
		}
		if postCreatePublishAt != "" {
			body["publishAt"] = postCreatePublishAt
		}
		if postCreateStatus != "" {
			body["status"] = postCreateStatus
		}
		if postCreateLabelIDs != "" {
			body["labelIds"] = splitCSV(postCreateLabelIDs)
		}
		if postCreateCampaignID != "" {
			body["campaignId"] = postCreateCampaignID
		}
		if postCreateNotes != "" {
			body["notes"] = postCreateNotes
		}
		// Validate required fields before authenticating / dialing the API.
		// The Ordinal API rejects posts missing any of these, and surfacing
		// the error locally gives a more actionable message without
		// consuming a rate-limit slot. Use a typed check so null, non-string,
		// and whitespace-only values fail locally too.
		for _, field := range []struct{ key, flag string }{
			{"title", "title"},
			{"publishAt", "publish-at"},
			{"status", "status"},
		} {
			s, ok := body[field.key].(string)
			if !ok || strings.TrimSpace(s) == "" {
				return fmt.Errorf("--%s or a non-empty %q field in --body-json/--body-file is required", field.flag, field.key)
			}
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewPostService(c).Create(body)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var postUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		body, err := parseBodyJSON(postUpdateBodyJSON, postUpdateBodyFile)
		if err != nil {
			return err
		}
		if body == nil {
			body = map[string]interface{}{}
		}
		if cmd.Flags().Changed("title") {
			body["title"] = postUpdateTitle
		}
		if cmd.Flags().Changed("publish-at") {
			body["publishAt"] = postUpdatePublishAt
		}
		if cmd.Flags().Changed("status") {
			body["status"] = postUpdateStatus
		}
		if cmd.Flags().Changed("label-ids") {
			body["labelIds"] = splitCSV(postUpdateLabelIDs)
		}
		if cmd.Flags().Changed("campaign-id") {
			body["campaignId"] = postUpdateCampaignID
		}
		if cmd.Flags().Changed("notes") {
			body["notes"] = postUpdateNotes
		}
		if len(body) == 0 {
			return fmt.Errorf("no fields to update; provide flags or --body-json/--body-file")
		}
		data, err := api.NewPostService(c).Update(postID, body)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var postArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a post (30-day deletion window)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewPostService(c).Archive(postID)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var postUnarchiveCmd = &cobra.Command{
	Use:   "unarchive",
	Short: "Restore a post from trash",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewPostService(c).Unarchive(postID)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var postScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule or reschedule a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewPostService(c).Schedule(postID, models.SchedulePostRequest{PublishAt: postSchedulePublishAt})
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var postUnscheduleCmd = &cobra.Command{
	Use:   "unschedule",
	Short: "Cancel a scheduled post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewPostService(c).Unschedule(postID)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}
