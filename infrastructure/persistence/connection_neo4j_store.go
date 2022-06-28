package persistence

import (
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const (
	FOLLOW_CONNECTION   = "FOLLOWS"
	USER_SKILL_RELATION = "HAS_SKILL"
)

type ConnectionNeo4jStore struct {
	Driver neo4j.Driver
}

func NewConnectionNeo4jStore(driver neo4j.Driver) ConnectionNeo4jStore {
	return ConnectionNeo4jStore{
		Driver: driver,
	}
}

func (u *ConnectionNeo4jStore) CreateConnectionBetweenUsers(toUser *domain.RegisteredUser, fromUser *domain.RegisteredUser) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.persistUserAsNode(tx, toUser)
	}); err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistUserAsNode(tx, fromUser)
		}); err != nil {
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistConnectionBetweenUsers(tx, fromUser, toUser)
		}); err != nil {
		return err
	}

	return nil
}

func (u *ConnectionNeo4jStore) AddSkillToUser(user *domain.RegisteredUser, skill string) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.persistUserAsNode(tx, user)
	}); err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistSkillAsNode(tx, skill)
		}); err != nil {
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistConnectionBetweenUserAndSkill(tx, user, skill)
		}); err != nil {
		return err
	}

	return nil
}

func (u *ConnectionNeo4jStore) AddJobOfferFromCompany(company *domain.Company, jobOffer *domain.JobOffer) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.persistCompanyAsNode(tx, company)
	}); err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistJobOfferAsNode(tx, jobOffer)
		}); err != nil {
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistConnectionBetweenCompanyAndJobOffer(tx, company, jobOffer)
		}); err != nil {
		return err
	}

	return nil
}

func (u *ConnectionNeo4jStore) AddRequiredSkillToJobOffer(skill string, jobOffer *domain.JobOffer) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.persistSkillAsNode(tx, skill)
	}); err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistJobOfferAsNode(tx, jobOffer)
		}); err != nil {
		return err
	}

	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistConnectionBetweenRequiredSkillAndJobOffer(tx, skill, jobOffer)
		}); err != nil {
		return err
	}

	return nil
}

func (u *ConnectionNeo4jStore) FindUsersConnection(username string) (connections []string, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = session.Close()
	}()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.findConnectionsByUsername(tx, username)
	})
	if result == nil {
		return nil, err
	}
	connections = result.([]string)
	return connections, err
}

func (u *ConnectionNeo4jStore) FindSuggestedConnectionsForUser(username string) (suggestions []string, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = session.Close()
	}()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.findSuggestedConnectionsForUser(tx, username)
	})
	if result == nil {
		return nil, err
	}
	suggestions = result.([]string)
	return suggestions, err
}

func (u *ConnectionNeo4jStore) DeleteConnection(usernameFrom string, usernameTo string) (ret interface{}, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = session.Close()
	}()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (:RegisteredUserNode {username: $usernameFrom})-[r:FOLLOWS]->(:RegisteredUserNode {username: $usernameTo}) DELETE r"
		parameters := map[string]interface{}{
			"usernameFrom": usernameFrom,
			"usernameTo":   usernameTo,
		}
		_, err := tx.Run(query, parameters)
		return nil, err
	})

	return nil, err
}

func (u *ConnectionNeo4jStore) persistUserAsNode(tx neo4j.Transaction, user *domain.RegisteredUser) (interface{}, error) {
	query := "MERGE (:RegisteredUserNode {email: $email, username: $username})"
	parameters := map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) persistSkillAsNode(tx neo4j.Transaction, skill string) (interface{}, error) {
	query := "MERGE (:SkillNode {name: $name})"
	parameters := map[string]interface{}{
		"name": skill,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) persistCompanyAsNode(tx neo4j.Transaction, company *domain.Company) (interface{}, error) {
	query := "MERGE (:CompanyNode {name: $name, username: $username, industry: $industry})"
	parameters := map[string]interface{}{
		"name":     company.CompanyName,
		"username": company.Username,
		"industry": company.Industry,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) persistJobOfferAsNode(tx neo4j.Transaction, offer *domain.JobOffer) (interface{}, error) {
	query := "MERGE (:JobOfferNode {position: $position, company: $company})"
	parameters := map[string]interface{}{
		"position": offer.Position,
		"company":  offer.Company.Username,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) persistConnectionBetweenUsers(tx neo4j.Transaction, fromUser *domain.RegisteredUser, toUser *domain.RegisteredUser) (interface{}, error) {
	query := "MATCH (from:RegisteredUserNode), (to:RegisteredUserNode) WHERE from.username = $fromUsername AND to.username = $toUsername CREATE (from)-[r:FOLLOWS]->(to)"
	parameters := map[string]interface{}{
		"fromUsername": fromUser.Username,
		"toUsername":   toUser.Username,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) persistConnectionBetweenUserAndSkill(tx neo4j.Transaction, user *domain.RegisteredUser, skill string) (interface{}, error) {
	query := "MATCH (u:RegisteredUserNode), (s:SkillNode) WHERE u.username = $user AND s.name = $skill CREATE (u)-[r:HAS_SKILL]->(s)"
	parameters := map[string]interface{}{
		"user":  user.Username,
		"skill": skill,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) persistConnectionBetweenCompanyAndJobOffer(tx neo4j.Transaction, company *domain.Company, offer *domain.JobOffer) (interface{}, error) {
	query := "MATCH (c:CompanyNode), (j:JobOfferNode) WHERE c.username = $company AND j.position = $position AND j.company = $company CREATE (c)-[r:OFFERS_JOB]->(j)"
	parameters := map[string]interface{}{
		"company":  company.Username,
		"position": offer.Position,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) persistConnectionBetweenRequiredSkillAndJobOffer(tx neo4j.Transaction, skill string, offer *domain.JobOffer) (interface{}, error) {
	query := "MATCH (s:SkillNode), (j:JobOfferNode) WHERE s.name = $skill AND j.position = $position AND j.company = $company CREATE (c)-[r:REQUIRES_SKILL]->(j)"
	parameters := map[string]interface{}{
		"skill":    skill,
		"position": offer.Position,
		"company":  offer.Company.Username,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *ConnectionNeo4jStore) findConnectionsByUsername(tx neo4j.Transaction, username string) ([]string, error) {
	records, err := tx.Run(
		"MATCH (u:RegisteredUserNode {username: $username})-[:FOLLOWS]->(connection) RETURN connection.username as usernameRet",
		map[string]interface{}{
			"username": username,
		},
	)
	if err != nil {
		return nil, err
	}

	var results []string
	for records.Next() {
		record := records.Record()
		username2, _ := record.Get("usernameRet")
		results = append(results, username2.(string))
	}

	return results, nil
}

func (u *ConnectionNeo4jStore) findSuggestedConnectionsForUser(tx neo4j.Transaction, username string) ([]string, error) {
	records, err := tx.Run(
		"MATCH (u:RegisteredUserNode {username: $username})-[r1:FOLLOWS]->(connection)-[r2:FOLLOWS]->(connection_of_connection) WHERE NOT connection_of_connection.username = $username RETURN DISTINCT connection_of_connection.username as usernameRet",
		map[string]interface{}{
			"username": username,
		},
	)
	if err != nil {
		return nil, err
	}

	var results []string
	for records.Next() {
		record := records.Record()
		username2, _ := record.Get("usernameRet")
		results = append(results, username2.(string))
	}

	return results, nil
}
