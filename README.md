rip (rest in peace)
===================

Allows to define REST API endpoints using a fluent nested interface:    
    
    api := rip.NewRIP()
    api.Path("items").Do(func(api *rip.Route) {
        api.POST().Target(map[string]interface{}{}).Handler(addItemHandler, "Add an item")
        api.GET().Param("ascending", "Order items in ascending order if true").
            Param("tag", "Only list items with given tags").
            Handler(listItemHandler, "List Items")
        api.Var("id").Do(func(api *rip.Route) {
            api.Delete().Handler(deleteItemHandler, "Remove an item")
        })
    })

Will create a http.Handler which allows adding items to the endpoint `/items`.
When calling `POST /items` the sent body (JSON) will be parsed into a
`map[string]interface{}`, defined by `Target()`. In the handler it can be retrieved like
`rip.Body(request)`.

It will also allow a GET `/items` which will return a list of items.
The query parameters `ascending` and `tag` cn be given to refine the result.
From the handlers the variables can be read by `rip.Vars(request).MustGetVar("tag")`.

Additionally an Item can be deleted by using a `DELETE /items/item-id` request.
The defined `id` var in the path can also be requested using
`rip.Vars(request).MustGetVar("id")` in the handler. 


