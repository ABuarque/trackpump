# trackpump 

![Tests](https://github.com/ABuarque/trackpump/workflows/Tests/badge.svg) ![Codecov](https://codecov.io/gh/ABuarque/trackpump/branch/master/graph/badge.svg?token=Y70u29hwMi) ![Deploy to production](https://github.com/ABuarque/trackpump/workflows/Deploy%20to%20production/badge.svg)

## Objective
Create a system to measure workout performance of a bodybuilder as weeks goes by.

## Infrastructure
One of the functional requirements of project is to run without costs. Getting the estimates of entities size we’ve got User having ~640 bytes and BodyMesurement ~364 bytes.  Looking for some database solutions available we have chosen [Google Datastore](https://cloud.google.com/datastore/pricing?hl=pt-br) as our database host due to the fact it provides to us 1GB free, which is a great space for us due to a year has 52 weeks and we have one BodyMesurement for week, thus 364 bytes x 52 weeks = 18928 bytes = 18,928 KB.
The the platform for our server is going to be [Google App Engine (GAE)](https://cloud.google.com/appengine/pricing) due to the fact it provides a free quota for its standard environment. GAE also provides scheduling solutions, and other functional requirement is to send automatic weekly reports about the performance. 
Other functional requirement is to keep pictures of body for each week, and the storage solution we’ve chosen was [pCloud](https://www.pcloud.com/pt/help/web-help-center/how-can-i-get-more-free-space) due to its free tier of 10 gb. 

## Metrics
- Weight (g);
- Abdominal Circumference (cm);
- Arm (cm);
- Forearm (cm);
- Calf (cm);
- Neck (cm);
- Hip (cm);
- Thigh (cm);

Besides these metrics two pictures of user should be taken to give a visual impression of progress. 

## Weekly Reports
Assuming the fact that user will not workout on sundays on that day user should get its weekly report by email. This report should show to user its progress on the collected metrics by showing to how its body fat percentage is, its body mass index and saying to him insights about those values: if it is necessary to lose or gain weight, for example. 
