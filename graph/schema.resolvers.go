package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/sql-judge/api-graphql/graph/generated"
	"github.com/sql-judge/api-graphql/graph/model"
)

func (r *mutationResolver) SubmitSolution(ctx context.Context, problemID int, solution string) (*model.Submission, error) {
	db := r.Resolver.db

	querySubmitSolution := `
		insert into submission (user_account_id, problem_id, solution)
		values (1, $1, $2)
		returning id
	`
	submissionRow := db.QueryRow(ctx, querySubmitSolution, problemID, solution)
	submission := &model.Submission{}
	if err := submissionRow.Scan(&submission.ID); err != nil {
		return nil, err
	}

	return submission, nil
}

func (r *queryResolver) Problem(ctx context.Context, id int) (*model.Problem, error) {
	db := r.Resolver.db

	problem := &model.Problem{}

	queryProblemInfo := `
		select id, title, description 
		from judge.problem 
		where id = $1`
	err := db.QueryRow(ctx, queryProblemInfo, id).Scan(&problem.ID, &problem.Title, &problem.Description)
	if err != nil {
		return nil, err
	}

	queryProblemAuthors := `
		select ua.id, ua.username, ua.full_name
		from problem p
		join problem_author pa on p.id = pa.problem_id
		join user_account ua on pa.user_account_id = ua.id
		where p.id = $1
		order by full_name nulls last`
	problemRows, err := db.Query(ctx, queryProblemAuthors, id)
	if err != nil {
		return nil, err
	}

	var authors []*model.User
	for problemRows.Next() {
		author := &model.User{}
		err := problemRows.Scan(&author.ID, &author.Username, &author.FullName)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	problem.Authors = authors

	queryProblemsTags := `
		select t.name, c.hex
		from problem p
		join problem_tag pt on p.id = pt.problem_id
		join tag t on pt.tag_id = t.id
		join color c on t.color_id = c.id
		where p.id = $1`
	tagsRows, err := db.Query(ctx, queryProblemsTags, id)
	if err != nil {
		return nil, err
	}

	for tagsRows.Next() {
		tag := &model.Tag{}
		err := tagsRows.Scan(&tag.Name, &tag.HexColor)
		if err != nil {
			return nil, err
		}

		problem.Tags = append(problem.Tags, tag)
	}

	return problem, nil
}

func (r *queryResolver) Problems(ctx context.Context) ([]*model.Problem, error) {
	db := r.Resolver.db

	// fetch all problems
	queryProblems := `
		select p.id, p.title, ps.accepted_ratio
		from problem p
		join problem_statistics ps on p.id = ps.id
		order by p.id`
	problemsRows, err := db.Query(ctx, queryProblems)
	if err != nil {
		return nil, err
	}

	var problems []*model.Problem
	for problemsRows.Next() {
		problem := &model.Problem{}

		err := problemsRows.Scan(&problem.ID, &problem.Title, &problem.AcceptedRatio)
		if err != nil {
			return nil, err
		}

		problems = append(problems, problem)
	}

	// fetch tags for all problems then add them to problems
	queryProblemsTags := `
		select p.id, t.name, c.hex
		from problem p
		join problem_tag pt on p.id = pt.problem_id
		join tag t on pt.tag_id = t.id
		join color c on t.color_id = c.id`
	tagsRows, err := db.Query(ctx, queryProblemsTags)
	if err != nil {
		return nil, err
	}

	idToProblem := make(map[int]*model.Problem, len(problems))
	for _, p := range problems {
		idToProblem[p.ID] = p
	}

	for tagsRows.Next() {
		var problemID int
		tag := &model.Tag{}
		err := tagsRows.Scan(&problemID, &tag.Name, &tag.HexColor)
		if err != nil {
			return nil, err
		}
		idToProblem[problemID].Tags = append(idToProblem[problemID].Tags, tag)
	}

	return problems, nil
}

func (r *queryResolver) Submissions(ctx context.Context) ([]*model.Submission, error) {
	db := r.Resolver.db

	// fetch all submissions
	querySubmissions := `
		select s.id, to_char(s.created_at, 'DD.MM.YYYY HH24:MI:SS:MS'), st.description, s.checker_message
		from submission s
		join submission_status st on s.status_id = st.id
		order by s.created_at desc, s.id desc
	`
	submissionsRows, err := db.Query(ctx, querySubmissions)
	if err != nil {
		return nil, err
	}

	var submissions []*model.Submission
	for submissionsRows.Next() {
		submission := &model.Submission{}
		err := submissionsRows.Scan(&submission.ID, &submission.CreatedAt, &submission.Status, &submission.CheckerMessage)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}

	// fetch problems info for all submissions
	querySubmissionsProblems := `
		select s.id, s.problem_id, p.title
		from submission s
		join problem p on s.problem_id = p.id
	`
	problemsRows, err := db.Query(ctx, querySubmissionsProblems)
	if err != nil {
		return nil, err
	}

	idToSubmission := make(map[int]*model.Submission, len(submissions))
	for _, s := range submissions {
		idToSubmission[s.ID] = s
	}

	for problemsRows.Next() {
		var submissionID int
		problem := &model.Problem{}
		err := problemsRows.Scan(&submissionID, &problem.ID, &problem.Title)
		if err != nil {
			return nil, err
		}
		idToSubmission[submissionID].Problem = problem
	}

	return submissions, nil
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	sql := "select id, username, full_name from user_account where id = $1"
	err := r.Resolver.db.QueryRow(ctx, sql, id).Scan(&user.ID, &user.Username, &user.FullName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
