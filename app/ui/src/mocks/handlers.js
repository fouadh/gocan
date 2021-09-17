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
    })

]

