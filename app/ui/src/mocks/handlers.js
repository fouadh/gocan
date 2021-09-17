import {rest} from 'msw'

export const handlers = [
    rest.get('/api/scenes', (req, res, ctx) => {
        return res(
            ctx.json({
                scenes: [
                    {
                        id: "s123",
                        name: "scene 1"
                    }]
            })
        )
    }),

    rest.get('/api/scenes/s123', (req, res, ctx) => {
        return res(ctx.json({
            id: "s123",
            name: "scene 1",
            applications: [
                {
                    id: "a1"
                }
            ]
        }))
    }),

    rest.get('/api/scenes/s123/apps', (req, res, ctx) => {
        return res(
            ctx.json({
                apps: [
                    {
                        id: "a1",
                        name: "app 1",
                        numberOfCommits: 123,
                        numberOfEntities: 45,
                        numberOfEntitiesChanged: 67,
                        numberOfAuthors: 89
                    }
                ]
            })
        )
    }),

    rest.get('/api/scenes/s123/apps/a1', (req, res, ctx) => {
        return res(ctx.json({
            name: "app 1"
        }))
    }),

    rest.get('/api/scenes/s123/apps/a1/revisions', (req, res, ctx) => {
        return res(ctx.json({
            revisions: [
                {entity: "file 1", numberOfRevisions: 10},
                {entity: "file 2", numberOfRevisions: 10},
                {entity: "file 3", numberOfRevisions: 8},
                {entity: "file 4", numberOfRevisions: 7},
                {entity: "file 5", numberOfRevisions: 3},
                {entity: "file 6", numberOfRevisions: 3},
                {entity: "file 7", numberOfRevisions: 2},
                {entity: "file 8", numberOfRevisions: 1},
                {entity: "file 9", numberOfRevisions: 1},
                {entity: "file 10", numberOfRevisions: 1},
                {entity: "file 11", numberOfRevisions: 1},
            ]
        }))
    }),

    rest.get('/api/scenes/s123/apps/a1/code-churn', (req, res, ctx) => {
        return res(ctx.json({
            codeChurn: [
                {date: "2021-09-01", added: 10, deleted: 5},
                {date: "2021-09-02", added: 100, deleted: 15},
                {date: "2021-09-03", added: 37, deleted: 51},
                {date: "2021-09-04", added: 48, deleted: 18},
                {date: "2021-09-05", added: 89, deleted: 0},
                {date: "2021-09-08", added: 250, deleted: 19},
                {date: "2021-09-10", added: 101, deleted: 45},
                {date: "2021-09-11", added: 36, deleted: 14},
                {date: "2021-09-12", added: 8, deleted: 63},
                {date: "2021-09-13", added: 123, deleted: 25},
                {date: "2021-09-14", added: 89, deleted: 3},
                {date: "2021-09-15", added: 21, deleted: 14},
                {date: "2021-09-16", added: 26, deleted: 41},
            ]
        }))
    })

]

