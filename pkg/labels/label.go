package labels

const (
	// CommunityContribution indicates that the contribution to the project
	// is not from a recognized party stated in the OWNERS or OWNERS_ALIASES file.
	CommunityContribution = "community-contribution"

	// InternalContribution indicates that the contribution to the project
	// is from a recognized party stated in the OWNERS or OWNERS_ALIASES file.
	InternalContribution = "internal-contribution"

	// Triaged indicates that the PR has been viewed and discussed
	// by the team to ensure that this is a legitimate change to the project.
	Triaged = "triaged"

	// Untriaged indicates that the PR has not been viewed and discussed
	// by the team to ensure that this is a legitimate change to the project.
	Untriaged = "untriaged"
)
