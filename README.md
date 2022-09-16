# aislack
A Slack bot that can generate AI art using the [Stable Diffusion](https://github.com/CompVis/stable-diffusion) model. Might also do other AI stuff in the future.

It passes image requests to [jobmgr](https://github.com/thatoddmailbox/jobmgr), so the bot and the actual image generation can run on separate machines. (and so that you can share the GPU with other programs)

## Things to do
* Error handling
* Randomize seed
* Make jobmgr endpoint configurable
* Advanced options (seed selection, CFG, step count, etc.) - probably want to make a Slack modal?