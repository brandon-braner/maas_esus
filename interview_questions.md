# Final Writeup

Hello,

I hope this message finds you well. I wanted to take a moment to express my sincere gratitude for the opportunity to work on this take-home assignment. I thoroughly enjoyed diving into the challenge and found it to be both engaging and intellectually stimulating. The problem space was fascinating, and I appreciated the chance to apply my skills in Go to something other then a Leet-code stye interview.

I have provided detailed responses to questions 3 and 4, outlining my approach to scaling the API service and integrating the new Memes AI feature. I look forward to your feedback and the opportunity to discuss my solutions further.

I want to take a quick minute to explain why I used the technologies I choose to use. 

### API

For the api I choose to use the build in net/http package in Go. Since 1.22 it handles 99% of what a basic api needs to handle. A more feature rich api framework may have made some other decisions about how to do things easier but I was concerned with speed / requests per minute in this exercise.

### Databases

**MongoDB**

For long term storage I choose to use MongoDB. I usually reach for Postgres but I saw this as an opportunity to do a bit more learning about using Mongo with Go. Also with the plan to use Redis to store tokens and other cache items and Mongos ability to handle transactions for things such as token updates in the database it felt like a good choice for scalability once we need to support 10,000 requests a second. 

**Redis / MemoryStore**

I wanted to use an in memory database to handle caching things like geolocation results as well as a fast way to lookup and charge tokens. I implemented the geolocation caching to get my requests per second to 100 as Googles reverse lookup takes between 250-300ms.

**Message Queue**

Once we need to scale we will need to be able to use event driven microservices to offload requests to other services and not wait for their response to continue our work. An example would be updating the tokens in redis in our meme service but putting a message on the queue to decrease or increase token counts. I would use AWS SQS for this to keep it on AWS for network latency. It also supports FIFO ordering so things that require the messages to be read in order will be.

# Question 3

Explain in as much detail as you like how you would scale this API service to ultimately
support a volume of 10,000 requests per second. Some things to consider include:

- How will you handle CI/CD in the context of a live service?
- How will you model and support SLAs? What should operational SLAs be for this
service?
- How do you support geographically diverse clients? As you scale the system out
horizontally, how do you continue to keep track of tokens without slowing down
the system?

### **CI/CD**

As we get to the point of supporting 10,000 requests a second I am going to suggest we would be on a container orchestration platform like Kubernetes or using a serverless hosting platform that allows for autoscaling like Fargate.

I have more experience with K8’s so I will explain my process for CI/CD with that and performing rollouts while keeping the application live. I would use Kubernetes, Helm, ArgoCD, Github Actions and Terraform.

First I would run any linting and tests we need to run prior to kicking off a build via Github Actions. GitHub Actions then automates the build and push process: when changes are pushed to the repository, a workflow triggers to build the Docker image for the web app, tag it with the new version (e.g., Git commit SHA), and push it to a container registry like Docker Hub or Amazon ECR. The workflow can also update the Helm chart’s `values.yaml` file with the new image tag and commit the change to the repository. I would also have Github actions run the Terraform changes during the pipeline so our environment is ready with any changes we need.

Next, ArgoCD handles the deployment to the Kubernetes cluster. ArgoCD continuously monitors the repository for changes to the Helm chart. When it detects an update (e.g., the new image tag in `values.yaml`), it automatically applies the changes to the cluster using the Helm chart. This ensures that the new version of the web app is deployed seamlessly. ArgoCD also provides rollback capabilities and health checks, ensuring the deployment is stable. This combination of tools creates a fully automated, GitOps-driven pipeline for deploying web app updates.

### SlAs

I would define SLA’s for this service based on the business goals the service was created for as well as operational goals that our technology org has aligned on. We already have one defined with 100 requests a second so we know that someone somewhere in the business cares about that.

Some SLA’s we could look at would include:

- **Availability** - how many 9’s are we looking for with this service. Keeping in mind the complexity each extra 9 add’s and we can’t go above the availability of our external providers I would set this somewhere around 99.99%.  We would monitor this through services like Cloudwatch or Promethus and set up alerts to make sure we stay in complience.
- **Throughput**: We have already discussed throughput with the 100 requests per second, however as our customer base grows we will want to increase this number. Example being the 10,000 per second we need to handle. This can also be measure with Cloudwatch or Promethus.
- **Error Rate -** We should probably keep an eye on error rates as well. If we start seeing a spike in errors it could mean either our code or one of our dependencies is having troubles and we need to investigate quickly to restore service. This also means we need to be logging errors and returning proper response codes to track these items.
- Message queue throughput and dead lettering. We need to make sure we are consuming our messages fast enough to keep our source of truth(Mongodb) up to date so we don’t lose track of user tokens if we lose them in Redis.
- **Developer Best Practices -** We want to make sure the developers are following best practices especially when it comes to things like code reviews and testing for these SLA’s.

Other things to think about is how to balance our SLA’s so they aren’t just vanity metrics. For example we could easily scale to 10,000 requests per second if we remove the service like Googles geocoding service that is taking 500ms to respond, but that isn’t the point of the sla. This is where the other sla like proper code reviews comes in.

## Scaling Across Regions

As we are using AWS cloud for this scaling across regions does become a bit easier but not without its challenges.  

We are using Kubernetes so we can deploy an EKS cluster to each region. We can do this by using global load balancing with Route 53 and application load balancers. 

If we stick with Mongodb we could use Mongodb atlas to handle multi region scaling. We could also move to a service like Dynamodb and global tables if we stay with nosql.

Memorystore does become a bit tricker as you can’t create multiregion instances of it. This makes it very important we keep our source of truth database up to date for things like tokens incase someone is traveling and they go into another region and use a different cluster.

# Meme’s AI

GenAI everywhere 🙂. I would handle this with roles / permissions through whoever our identity provider is assuming they offer authorization on top of authentication. Customers who sign up for the ai plan will get roles put onto their user and subsequently their jwt’s. We can inspect the jwt’s during their request to see if they have the role to use the gen ai meme service. I would also suspect that maybe these memes use more tokens so I would update the app to know how many tokens to charge based on the type of meme it generates.

On a side note I actually started this project by using OpenAI to generate the memes so I had to work backward a bit on this one. In the repo you will see there is a generators.go file. This has a interface in it called `MemeGenerator` that has a method `generate`. I then have two structs `TextMemeGenerator` and `AITextMemeGenerator` that are responsible for generating memes with and without ai. Based on the users permission that we get from the jwt in the auth middleware and pass into the request context it will generate the appropriate type of meme for them. It currently only supports OpenAI but could implement a strategy pattern to allow us to support multiple.

I look forward to talking with you all further about this project. 

Thank you once again for this opportunity.

Best regards,
Brandon Braner
brandon.braner@gmail.com
952-452-6872