/*
	This example shows the usage of a Paginator interface
*/
package main

import (
	"errors"
	"fmt"

	"github.com/exactlylabs/go-rest/pkg/restapi"
	"github.com/exactlylabs/go-rest/pkg/restapi/paginator"
	"github.com/exactlylabs/go-rest/pkg/restapi/webcontext"
)

// In our fake DB, this is the number of rows that we have there
const RESOURCES_NUMBER = 20

func main() {
	server, err := restapi.NewWebServer()
	if err != nil {
		panic(err)
	}
	server.Route("/", List, "GET")
	server.Run("127.0.0.1:5000")
}

// This is the resource returned in the collection
type Resource struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ResourceIterator struct {
	pointer int
	items   []Resource
}

func (i *ResourceIterator) Count() (uint64, error) {
	// Count is not the number of items in the iterator, but rather the total number of resources that are stored in the DB
	return RESOURCES_NUMBER, nil
}

func (i *ResourceIterator) GetRow() (Resource, error) {
	if !i.HasNext() {
		// Shouldn't be called anymore
		return Resource{}, errors.New("EOF")
	}
	resource := i.items[i.pointer]
	i.pointer++
	return resource, nil
}

func (i *ResourceIterator) HasNext() bool {
	return i.pointer < len(i.items)
}

func GetResourceCollection(limit, offset int) *ResourceIterator {
	items := make([]Resource, 0)
	for i := offset; i < offset+limit && i < RESOURCES_NUMBER; i++ {
		items = append(items, Resource{Id: i, Name: fmt.Sprintf("Item %d", i)})
	}
	return &ResourceIterator{items: items}
}

// List simulates a normal endpoint where someone wants to query a list of resources from a collection
// Here, we added a pagination functionality. For it to work, you need to instantiate a new Paginator object, telling it what is the type of the resource to be paginated
// You can customize the Default Limit and Offset, as well as the maximum Limit value
// Then, just call Paginate, providing a function that returns an iterator of the Resource type, based on the given limit and offset arguments
// You can see that we also included a "paginatorArgs" argument. That's because there's a Dependency provider built-in in the framework that parses those args for us.
// This argument is not required, it's here just to show that you can have access to the limit and offset values provided by the user.
func List(ctx *webcontext.Context, paginatorArgs *paginator.PaginationArgs) {
	// Play with this by changing the "limit" and "offset" query parameters
	p := paginator.New[Resource]()
	p.MaxLimit = 5
	paginatedResp, err := p.Paginate(ctx, func(limit, offset int) (paginator.Iterator[Resource], error) {
		return GetResourceCollection(limit, offset), nil
	})
	if err == nil && ctx.HasErrors() {
		return
	}

	ctx.JSON(200, paginatedResp)
}
