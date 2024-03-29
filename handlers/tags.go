package handlers

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"unicode"

	"cloud.google.com/go/firestore"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"google.golang.org/api/iterator"
)

func AddUserTags(ctx context.Context, ids []*model.UserDetails, tags []*string) (*bool, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)

	for _, vvv := range ids {
		value := vvv
		if value == nil {
			continue
		}
		userLspID := value.UserLspID
		userId := value.UserID
		if tags == nil || userLspID == nil || userId == nil {
			return nil, errors.New("please enter all the values of userLspId, userId and tags")
		}
		id := *userLspID
		uId := *userId
		var tagsArray []string
		for _, vv := range tags {
			if vv == nil {
				continue
			}
			v := *vv
			v = strings.ToLower(v)
			if isASCII(v) {
				tagsArray = append(tagsArray, v)
			} else {
				return nil, errors.New("please enter only ASCII values in tags")
			}
		}
		_, err = global.Client.Collection("userLspIdTags").Doc(id).Set(ctx, map[string]interface{}{
			"Tags":   tagsArray,
			"UserId": uId,
			"LspId":  lspId,
		})
		if err != nil {
			return nil, err
		}

	}

	res := true
	return &res, nil
}

func GetUserLspIDTags(ctx context.Context, userLspID []*string) ([]*model.TagsData, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if len(userLspID) == 0 {
		return nil, errors.New("please provide userLspIds")
	}

	res := make([]*model.TagsData, len(userLspID))
	var wg sync.WaitGroup
	for kk, vvv := range userLspID {
		vv := vvv
		if vv == nil {
			continue
		}
		wg.Add(1)
		go func(v *string, k int) {
			defer wg.Done()
			snap, err := global.Client.Collection("userLspIdTags").Doc(*v).Get(ctx)
			if err != nil {
				log.Printf("Got error while getting data: %v", err)
				return
			}
			data := snap.Data()
			tags := data["Tags"].([]interface{})
			userId := data["UserId"].(string)
			lspId := data["LspId"].(string)

			var tagsArray []*string
			for _, v := range tags {
				tmp := v.(string)
				tagsArray = append(tagsArray, &tmp)
			}

			tmp := model.TagsData{
				UserLspID: v,
				UserID:    &userId,
				Tags:      tagsArray,
				LspID:     &lspId,
			}
			res[k] = &tmp

		}(vv, kk)
	}
	wg.Wait()

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

func GetTagUsers(ctx context.Context, prevPageSnapShot *string, pageSize *int, tags []*string) (*model.PaginatedTagsData, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	lspId := claims["lsp_id"].(string)
	if len(tags) == 0 {
		return nil, errors.New("please provide tags")
	}
	if pageSize == nil {
		*pageSize = 10
	}
	var tmp []string
	var tagsArray []string
	for _, vv := range tags {
		if vv == nil {
			continue
		}
		v := *vv
		v = strings.ToLower(v)
		tagsArray = append(tagsArray, v)
	}

	var iter *firestore.DocumentIterator
	//prevPageSnapShot.(firestore.DocumentRef)
	if prevPageSnapShot != nil {
		iter = global.Client.Collection("userLspIdTags").Where("Tags", "array-contains-any", tagsArray).Limit(*pageSize).StartAfter(prevPageSnapShot).Where("LspId", "==", lspId).OrderBy("UserId", firestore.Asc).Documents(ctx)
	} else {
		iter = global.Client.Collection("userLspIdTags").Where("Tags", "array-contains-any", tagsArray).Limit(*pageSize).Where("LspId", "==", lspId).OrderBy("__name__", firestore.Asc).Documents(ctx)
	}

	var maps []map[string]interface{}

	var snapShot string
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

		snapShot = doc.Data()["UserId"].(string)
		tmp = append(tmp, doc.Ref.ID)
		maps = append(maps, doc.Data())
	}
	if maps == nil {
		return nil, nil
	}
	if len(maps) == 0 {
		return nil, nil
	}

	res := make([]*model.TagsData, len(maps))
	var wg sync.WaitGroup
	for kk, vvv := range maps {
		wg.Add(1)
		vv := vvv
		go func(k int, v map[string]interface{}) {
			defer wg.Done()
			tagsInterface := v["Tags"].([]interface{})
			var allTags []*string
			for _, v := range tagsInterface {
				tmp := v.(string)
				allTags = append(allTags, &tmp)
			}

			userLspId := tmp[k]
			userId := v["UserId"].(string)
			tmp := model.TagsData{
				UserLspID: &userLspId,
				UserID:    &userId,
				Tags:      allTags,
				LspID:     &lspId,
			}
			res[k] = &tmp

		}(kk, vv)

	}
	wg.Wait()

	resp := model.PaginatedTagsData{
		Data:             res,
		PrevPageSnapShot: &snapShot,
	}
	return &resp, nil
}
