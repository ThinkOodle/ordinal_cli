package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	ideaID                      string
	ideaListLimit               int
	ideaListCursor              string
	ideaListIDs                 string
	ideaListChannel             string
	ideaListLinkedInProfileID   string
	ideaListXProfileID          string
	ideaListLabelIDs            string
	ideaListCreatedAtMin        string
	ideaListCreatedAtMax        string
	ideaListSortBy              string
	ideaListSortOrder           string
	ideaListAll                 bool
	ideaCreateTitle             string
	ideaCreateLabelIDs          string
	ideaCreateCampaignID        string
	ideaCreateBodyJSON          string
	ideaCreateBodyFile          string
	ideaUpdateTitle             string
	ideaUpdateLabelIDs          string
	ideaUpdateCampaignID        string
	ideaUpdateBodyJSON          string
	ideaUpdateBodyFile          string
	ideaAddToCalendarPublishDate string
)

func init() {
	rootCmd.AddCommand(ideaCmd)
	ideaCmd.AddCommand(ideaListCmd)
	ideaCmd.AddCommand(ideaGetCmd)
	ideaCmd.AddCommand(ideaCreateCmd)
	ideaCmd.AddCommand(ideaUpdateCmd)
	ideaCmd.AddCommand(ideaArchiveCmd)
	ideaCmd.AddCommand(ideaUnarchiveCmd)
	ideaCmd.AddCommand(ideaAddToCalendarCmd)

	ideaListCmd.Flags().IntVar(&ideaListLimit, "limit", 0, "Max ideas to return (1-100)")
	ideaListCmd.Flags().StringVar(&ideaListCursor, "cursor", "", "Pagination cursor")
	ideaListCmd.Flags().StringVar(&ideaListIDs, "ids", "", "Fetch specific ideas by UUIDs (comma-separated)")
	ideaListCmd.Flags().StringVar(&ideaListChannel, "channel", "", "Filter by channel")
	ideaListCmd.Flags().StringVar(&ideaListLinkedInProfileID, "linkedin-profile-id", "", "Filter by LinkedIn profile ID")
	ideaListCmd.Flags().StringVar(&ideaListXProfileID, "x-profile-id", "", "Filter by X/Twitter profile ID")
	ideaListCmd.Flags().StringVar(&ideaListLabelIDs, "label-ids", "", "Filter by label IDs (comma-separated)")
	ideaListCmd.Flags().StringVar(&ideaListCreatedAtMin, "created-at-min", "", "Filter ideas created on or after this date")
	ideaListCmd.Flags().StringVar(&ideaListCreatedAtMax, "created-at-max", "", "Filter ideas created on or before this date")
	ideaListCmd.Flags().StringVar(&ideaListSortBy, "sort-by", "", "Field to sort by")
	ideaListCmd.Flags().StringVar(&ideaListSortOrder, "sort-order", "", "Sort order (asc or desc)")
	ideaListCmd.Flags().BoolVar(&ideaListAll, "all", false, "Auto-paginate and return all results")

	ideaGetCmd.Flags().StringVar(&ideaID, "id", "", "Idea ID (UUID)")
	ideaGetCmd.MarkFlagRequired("id")

	ideaCreateCmd.Flags().StringVar(&ideaCreateTitle, "title", "", "Idea title")
	ideaCreateCmd.Flags().StringVar(&ideaCreateLabelIDs, "label-ids", "", "Comma-separated label IDs")
	ideaCreateCmd.Flags().StringVar(&ideaCreateCampaignID, "campaign-id", "", "Campaign ID")
	ideaCreateCmd.Flags().StringVar(&ideaCreateBodyJSON, "body-json", "", "Full JSON body including channel configs (individual flags override matching keys)")
	ideaCreateCmd.Flags().StringVar(&ideaCreateBodyFile, "body-file", "", "Path to JSON body file (or - for stdin)")

	ideaUpdateCmd.Flags().StringVar(&ideaID, "id", "", "Idea ID (UUID)")
	ideaUpdateCmd.Flags().StringVar(&ideaUpdateTitle, "title", "", "Idea title")
	ideaUpdateCmd.Flags().StringVar(&ideaUpdateLabelIDs, "label-ids", "", "Comma-separated label IDs")
	ideaUpdateCmd.Flags().StringVar(&ideaUpdateCampaignID, "campaign-id", "", "Campaign ID")
	ideaUpdateCmd.Flags().StringVar(&ideaUpdateBodyJSON, "body-json", "", "Full JSON body (individual flags override matching keys when set)")
	ideaUpdateCmd.Flags().StringVar(&ideaUpdateBodyFile, "body-file", "", "Path to JSON body file (or - for stdin)")
	ideaUpdateCmd.MarkFlagRequired("id")

	ideaArchiveCmd.Flags().StringVar(&ideaID, "id", "", "Idea ID (UUID)")
	ideaArchiveCmd.MarkFlagRequired("id")
	ideaUnarchiveCmd.Flags().StringVar(&ideaID, "id", "", "Idea ID (UUID)")
	ideaUnarchiveCmd.MarkFlagRequired("id")

	ideaAddToCalendarCmd.Flags().StringVar(&ideaID, "id", "", "Idea ID (UUID)")
	ideaAddToCalendarCmd.Flags().StringVar(&ideaAddToCalendarPublishDate, "publish-date", "", "Date to add to calendar (YYYY-MM-DD or ISO 8601)")
	ideaAddToCalendarCmd.MarkFlagRequired("id")
	ideaAddToCalendarCmd.MarkFlagRequired("publish-date")
}

var ideaCmd = &cobra.Command{
	Use:   "idea",
	Short: "Manage content ideas",
}

var ideaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ideas",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		params := models.ListIdeasParams{
			Limit:             ideaListLimit,
			Cursor:            ideaListCursor,
			IDs:               ideaListIDs,
			Channel:           ideaListChannel,
			LinkedInProfileID: ideaListLinkedInProfileID,
			XProfileID:        ideaListXProfileID,
			LabelIDs:          ideaListLabelIDs,
			CreatedAtMin:      ideaListCreatedAtMin,
			CreatedAtMax:      ideaListCreatedAtMax,
			SortBy:            ideaListSortBy,
			SortOrder:         ideaListSortOrder,
		}
		svc := api.NewIdeaService(c)
		if ideaListAll {
			ideas, err := svc.ListAll(params)
			if err != nil {
				return err
			}
			return printResult(ideas)
		}
		resp, err := svc.List(params)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var ideaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an idea by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		i, err := api.NewIdeaService(c).Get(ideaID)
		if err != nil {
			return err
		}
		return printResult(i)
	},
}

var ideaCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an idea",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		body, err := parseBodyJSON(ideaCreateBodyJSON, ideaCreateBodyFile)
		if err != nil {
			return err
		}
		if body == nil {
			body = map[string]interface{}{}
		}
		if ideaCreateTitle != "" {
			body["title"] = ideaCreateTitle
		}
		if ideaCreateLabelIDs != "" {
			body["labelIds"] = splitCSV(ideaCreateLabelIDs)
		}
		if ideaCreateCampaignID != "" {
			body["campaignId"] = ideaCreateCampaignID
		}
		if _, ok := body["title"]; !ok {
			return fmt.Errorf("--title or a title field in --body-json/--body-file is required")
		}
		data, err := api.NewIdeaService(c).Create(body)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var ideaUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an idea",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		body, err := parseBodyJSON(ideaUpdateBodyJSON, ideaUpdateBodyFile)
		if err != nil {
			return err
		}
		if body == nil {
			body = map[string]interface{}{}
		}
		if cmd.Flags().Changed("title") {
			body["title"] = ideaUpdateTitle
		}
		if cmd.Flags().Changed("label-ids") {
			body["labelIds"] = splitCSV(ideaUpdateLabelIDs)
		}
		if cmd.Flags().Changed("campaign-id") {
			body["campaignId"] = ideaUpdateCampaignID
		}
		if len(body) == 0 {
			return fmt.Errorf("no fields to update; provide flags or --body-json/--body-file")
		}
		data, err := api.NewIdeaService(c).Update(ideaID, body)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var ideaArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive an idea",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewIdeaService(c).Archive(ideaID)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var ideaUnarchiveCmd = &cobra.Command{
	Use:   "unarchive",
	Short: "Unarchive an idea",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewIdeaService(c).Unarchive(ideaID)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var ideaAddToCalendarCmd = &cobra.Command{
	Use:   "add-to-calendar",
	Short: "Convert an idea into a scheduled post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewIdeaService(c).AddToCalendar(ideaID, models.AddIdeaToCalendarRequest{PublishDate: ideaAddToCalendarPublishDate})
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}
