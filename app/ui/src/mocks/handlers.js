import { rest } from 'msw'

export const handlers = [
    rest.get('/api/scenes', (req, res, ctx) => {
        return res(
            ctx.json({
                "scenes": [
                    {
                        "id": "123",
                        "name": "scene 1",
                        "applications": []
                    }]
            })
        )
    }),
]

