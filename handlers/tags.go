package handlers

import (
	"context"
	"errors"
	"log"
	"unicode"

	"github.com/zicops/zicops-notification-server/global"
	"google.golang.org/api/iterator"
)

func AddUserTags(ctx context.Context, userLspID *string, tags []*string) (*bool, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if tags == nil || userLspID == nil {
		return nil, errors.New("please enter both userLspId and tags")
	}
	id := *userLspID
	var tagsArray []string
	for _, vv := range tags {
		v := *vv
		if isASCII(v) {
			tagsArray = append(tagsArray, v)
		} else {
			return nil, errors.New("please enter only ASCII values in tags")
		}
	}
	_, err = global.Client.Collection("userLspIdTags").Doc(id).Set(ctx, map[string]interface{}{
		"Tags": tagsArray,
	})
	if err != nil {
		return nil, err
	}
	res := true
	return &res, nil
}

func GetUserLspIDTags(ctx context.Context, userLspID *string) ([]*string, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	snap, err := global.Client.Collection("userLspIdTags").Doc(*userLspID).Get(ctx)
	if err != nil {
		return nil, err
	}
	data := snap.Data()
	tags := data["Tags"].([]interface{})

	var res []*string
	for _, v := range tags {
		tmp := v.(string)
		res = append(res, &tmp)
	}
	return res, nil
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func GetTagUsers(ctx context.Context, tags []*string) ([]*string, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var tmp []string
	var tagsArray []string
	for _, vv := range tags {
		v := *vv
		tagsArray = append(tagsArray, v)
	}
	iter := global.Client.Collection("userLspIdTags").Where("Tags", "array-contains-any", tagsArray).Documents(ctx)

	for {
		doc, err := iter.Next()
		//see if iterator is done
		if err == iterator.Done {
			break
		}

		//see if the error is no more items in iterator
		if err != nil && err.Error() == "no more items in iterator" {
			break
			//return nil, nil
		}

		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return nil, err
		}
		tmp = append(tmp, doc.Ref.ID)
	}

	var res []*string
	for _, v := range tmp {
		vv := v
		res = append(res, &vv)
	}
	return res, nil
}
