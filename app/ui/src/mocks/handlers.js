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
    }),

    rest.get('/api/scenes/s123/apps/a1/boundaries', (req, res, ctx) => {
        return res(ctx.json({
            boundaries: [
                {id: "b123", name: "production/tests"}
            ]
        }))
    }),

    rest.get('/api/scenes/s123/apps/a1/revisions-trends?boundary=b123', (req, res, ctx) => {
        return res(ctx.json({
            trends: [
                {
                    date: "2021-09-10", revisions: [
                        {entity: "tests", numberOfRevisions: 50},
                        {entity: "production", numberOfRevisions: 100}
                    ]
                },{
                    date: "2021-09-11", revisions: [
                        {entity: "tests", numberOfRevisions: 35},
                        {entity: "production", numberOfRevisions: 78}
                    ]
                },{
                    date: "2021-09-12", revisions: [
                        {entity: "tests", numberOfRevisions: 250},
                        {entity: "production", numberOfRevisions: 210}
                    ]
                },{
                    date: "2021-09-13", revisions: [
                        {entity: "tests", numberOfRevisions: 198},
                        {entity: "production", numberOfRevisions: 86}
                    ]
                },{
                    date: "2021-09-14", revisions: [
                        {entity: "tests", numberOfRevisions: 74},
                        {entity: "production", numberOfRevisions: 210}
                    ]
                },{
                    date: "2021-09-15", revisions: [
                        {entity: "tests", numberOfRevisions: 36},
                        {entity: "production", numberOfRevisions: 12}
                    ]
                },{
                    date: "2021-09-16", revisions: [
                        {entity: "tests", numberOfRevisions: 145},
                        {entity: "production", numberOfRevisions: 432}
                    ]
                }
            ]
        }))
    }),

    rest.get('/api/scenes/s123/apps/a1/active-set', (req, res, ctx) => {
        return res(ctx.json({
            activeSet: [
                {date: '2021-09-10', opened: 123, closed: 362},
                {date: '2021-09-11', opened: 86, closed: 10},
                {date: '2021-09-12', opened: 89, closed: 50},
                {date: '2021-09-13', opened: 133, closed: 25},
                {date: '2021-09-14', opened: 75, closed: 63},
                {date: '2021-09-15', opened: 68, closed: 45},
                {date: '2021-09-16', opened: 85, closed: 47},
                {date: '2021-09-17', opened: 124, closed: 23},
            ]
        }))
    }),

    rest.get('/api/scenes/s123/apps/a1/developers', (req, res, ctx) => {
        return res(ctx.json({
            authors: [
                {name: "Alice"},
                {name: "Bob"},
                {name: "Charlie"},
            ]
        }))
    }),

    rest.get('/api/scenes/s123/apps/a1/knowledge-map', (req, res, ctx) => {
        return res(ctx.json({
                "name": "gocan",
                "children": [
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
                                                "weight": 1,
                                                "size": 45,
                                                "mainDeveloper": "Alice"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    },
                    {
                        "name": "business",
                        "children": [
                            {
                                "name": "core",
                                "children": [
                                    {
                                        "name": "coupling",
                                        "children": [
                                            {
                                                "name": "commands.go",
                                                "weight": 1,
                                                "size": 99,
                                                "mainDeveloper": "Bob"
                                            },
                                            {
                                                "name": "coupling.go",
                                                "weight": 1,
                                                "size": 43,
                                                "mainDeveloper": "Charlie"
                                            }
                                        ]
                                    },
                                    {
                                        "name": "developer",
                                        "children": [
                                            {
                                                "name": "commands.go",
                                                "weight": 1,
                                                "size": 49,
                                                "mainDeveloper": "Alice"
                                            },
                                            {
                                                "name": "developer.go",
                                                "weight": 1,
                                                "size": 32,
                                                "mainDeveloper": "Alice"
                                            }
                                        ]
                                    }
                                ]
                            },
                            {
                                "name": "data",
                                "children": [
                                    {
                                        "name": "schema",
                                        "children": [
                                            {
                                                "name": "folder",
                                                "children": [
                                                    {
                                                        "name": "file1",
                                                        "weight": 1,
                                                        "size": 1,
                                                        "mainDeveloper": "Alice"
                                                    },
                                                    {
                                                        "name": "file2",
                                                        "weight": 1,
                                                        "size": 24,
                                                        "mainDeveloper": "Alice"
                                                    },
                                                    {
                                                        "name": "file3",
                                                        "weight": 1,
                                                        "size": 1,
                                                        "mainDeveloper": "Bob"
                                                    },
                                                    {
                                                        "name": "file4",
                                                        "weight": 1,
                                                        "size": 41,
                                                        "mainDeveloper": "Charlie"
                                                    }
                                                ]
                                            }
                                        ]
                                    },
                                    {
                                        "name": "store",
                                        "children": [
                                            {
                                                "name": "another_folder",
                                                "children": [
                                                    {
                                                        "name": "file5",
                                                        "weight": 1,
                                                        "size": 79,
                                                        "mainDeveloper": "Charlie"
                                                    },
                                                    {
                                                        "name": "file6",
                                                        "weight": 1,
                                                        "size": 11,
                                                        "mainDeveloper": "Alice"
                                                    }
                                                ]
                                            },
                                            {
                                                "name": "test",
                                                "children": [
                                                    {
                                                        "name": "file7",
                                                        "weight": 1,
                                                        "size": 45,
                                                        "mainDeveloper": "Alice"
                                                    },
                                                    {
                                                        "name": "file8",
                                                        "weight": 1,
                                                        "size": 8,
                                                        "mainDeveloper": "Bob"
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
            }
        ))
    })
]
