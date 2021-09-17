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
    }),

    rest.get('/api/scenes/s123/apps/a1/coupling-hierarchy', (req, res, ctx) => {
        return res(ctx.json({
            "name": "root",
            "children": [
                {
                    "name": "business",
                    "children": [
                        {
                            "name": "core",
                            "children": [
                                {
                                    "name": "start-db",
                                    "children": [
                                        {
                                            "name": "command.go",
                                            "coupling": [
                                                "root/business/core/stop-db/command.go"
                                            ],
                                            "relations": [
                                                {
                                                    "coupled": "root/business/core/stop-db/command.go",
                                                    "degree": 0.8,
                                                    "averageRevisions": 5
                                                }
                                            ]
                                        }
                                    ]
                                },
                                {
                                    "name": "setup-db",
                                    "children": [
                                        {
                                            "name": "command.go",
                                            "coupling": [
                                                "root/business/core/start-db/command.go"
                                            ],
                                            "relations": [
                                                {
                                                    "coupled": "root/business/core/start-db/command.go",
                                                    "degree": 0.6,
                                                    "averageRevisions": 5
                                                }
                                            ]
                                        }
                                    ]
                                },
                                {
                                    "name": "create-scene",
                                    "children": [
                                        {
                                            "name": "command.go",
                                            "coupling": [
                                                "root/business/core/start-db/command.go"
                                            ],
                                            "relations": [
                                                {
                                                    "coupled": "root/business/core/start-db/command.go",
                                                    "degree": 0.6,
                                                    "averageRevisions": 5
                                                }
                                            ]
                                        },
                                        {
                                            "name": "usecase.go"
                                        }
                                    ]
                                },
                                {
                                    "name": "stop-db",
                                    "children": [
                                        {
                                            "name": "command.go"
                                        }
                                    ]
                                },
                                {
                                    "name": "coupling",
                                    "children": [
                                        {
                                            "name": "commands.go"
                                        }
                                    ]
                                }
                            ]
                        },
                        {
                            "name": "platform",
                            "children": [
                                {
                                    "name": "context.go"
                                }
                            ]
                        }
                    ]
                },
                {
                    "name": "app",
                    "children": [
                        {
                            "name": "cmd",
                            "children": [
                                {
                                    "name": "gocan",
                                    "children": [
                                        {
                                            "name": "main.go",
                                            "coupling": [
                                                "root/business/core/stop-db/command.go",
                                                "root/business/core/create-scene/command.go",
                                                "root/business/core/setup-db/command.go",
                                                "root/business/core/start-db/command.go",
                                                "root/business/platform/context.go",
                                                "root/business/core/coupling/commands.go",
                                                "root/business/core/create-scene/usecase.go"
                                            ],
                                            "relations": [
                                                {
                                                    "coupled": "root/business/core/stop-db/command.go",
                                                    "degree": 0.4444444444444444,
                                                    "averageRevisions": 9
                                                },
                                                {
                                                    "coupled": "root/business/core/create-scene/command.go",
                                                    "degree": 0.4444444444444444,
                                                    "averageRevisions": 9
                                                },
                                                {
                                                    "coupled": "root/business/core/setup-db/command.go",
                                                    "degree": 0.4444444444444444,
                                                    "averageRevisions": 9
                                                },
                                                {
                                                    "coupled": "root/business/core/start-db/command.go",
                                                    "degree": 0.4,
                                                    "averageRevisions": 10
                                                },
                                                {
                                                    "coupled": "root/business/platform/context.go",
                                                    "degree": 0.35294117647058826,
                                                    "averageRevisions": 8.5
                                                },
                                                {
                                                    "coupled": "root/business/core/coupling/commands.go",
                                                    "degree": 0.35294117647058826,
                                                    "averageRevisions": 8.5
                                                },
                                                {
                                                    "coupled": "root/business/core/create-scene/usecase.go",
                                                    "degree": 0.35294117647058826,
                                                    "averageRevisions": 8.5
                                                }
                                            ]
                                        }
                                    ]
                                }
                            ]
                        }
                    ]
                }
            ]
        }))
    }),

    rest.get('/api/scenes/s123/apps/a1/modus-operandi', (req, res, ctx) => {
        return res(ctx.json({
            modusOperandi: [
                {word: "the", count: 123},
                {word: "quick", count: 42},
                {word: "fox", count: 545},
                {word: "jumped", count: 236},
                {word: "over", count: 475},
                {word: "lazy", count: 368},
                {word: "dog", count: 742}
            ]
        }))
    })
]

