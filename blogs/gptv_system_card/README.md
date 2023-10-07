# GPT-4V(ision) System Card 论文研读

## 解读

OpenAI 放出了 19 页的 GPT-4V(ision)报告来解释这个模型，释放了大量信息，模型早在 2022 年就训练好了，现在才放出来主要是人工智能安全和合规考量：

1. GPT-4V 是 OpenAI 开发的一个具有视觉能力的语言模型,能够分析用户提供的图像输入并指示 GPT-4 进行分析。它结合了文本和视觉两种模式,拓展了仅限文本的系统的影响力和风险范围。

1. OpenAI 采取了渐进式部署方法,首先让一小部分用户试用,以获得反馈和洞察真实的交互方式。这有助于 OpenAI 认识到一些风险,如模型的误报和限制、人脸识别的隐私考量等。

1. OpenAI 进行了定性和定量评估以了解系统,包括聘请外部专家进行军事化测试,并建立了评估模型拒绝率和性能准确性的指标。重点评估领域包括有害内容、代表性、分配和服务质量的风险、隐私、网络安全、多模态越狱等。

1. 评估发现了科学、医学建议、刻板印象、没有根据的推断等方面的一些限制,需要采取缓解措施。

1. OpenAI 采取了模型级和系统级的缓解措施,通过额外的安全训练数据增强了对非法行为和无根据推断请求的拒绝行为,并增加了针对包含文字的对抗图像的系统级缓解措施。

下一步 OpenAI 将继续关注是否应允许模型进行某些行为、提高全球用户使用的语言和图像识别能力、获取更高精度的人像处理能力等方面。

## 原文

> address: https://cdn.openai.com/papers/GPTV_System_Card.pdf

### 1 Introduction

GPT-4 with vision (GPT-4V) enables users to instruct GPT-4 to analyze image inputs provided by the user, and is the latest capability we are making broadly available. Incorporating additional modalities (such as image inputs) into large language models (LLMs) is viewed by some as a key frontier in artificial intelligence research and development. Multimodal LLMs offer the possibility of expanding the impact of language-only systems with novel interfaces and capabilities, enabling them to solve new tasks and provide novel experiences for their users.

In this system card, we analyze the safety properties of GPT-4V. Our work on safety for GPT-4V builds on the work done for GPT-4 and here we dive deeper into the evaluations, preparation, and _mitigation_ work done specifically for image inputs.

Similar to GPT-4, training of GPT-4V was completed in 2022 and we began providing early access to the system in March 2023. As GPT-4 is the technology behind the visual capabilities of GPT-4V, its training process was the same. The pre-trained model was first trained to predict the next word in a document, using a large dataset of text and image data from the Internet as well as licensed sources of data. It was then fine-tuned with additional data, using an algorithm called reinforcement learning from human feedback (RLHF), to produce outputs that are preferred by human trainers.

Large multimodal models introduce different limitations and expand the risk surface compared to text-based language models. GPT-4V possesses the limitations and capabilities of each modality (text and vision), while at the same time presenting novel capabilities emerging from the intersection of said modalities and from the intelligence and reasoning afforded by large scale models.

This system card outlines how OpenAI prepared the vision capabilities of GPT-4 for deployment. It describes the early access period of the model for small scale users and safety learnings OpenAI gained from this period, multimodal evaluations built to study the model's fitness for deployment, key findings of expert red teamers, and the _mitigations_ OpenAI implemented prior to broad release.

### 2 Deployment Preparation

**2.1 Learnings from early access**

OpenAI gave a diverse set of alpha users access to GPT-4V earlier this year, including Be My Eyes, an organization that builds tools for visually impaired users.

_2.1.1 Be My Eyes_

Beginning in March, 2023, Be My Eyes and OpenAI collaborated to develop Be My AI, a new tool to describe the visual world for people who are blind or have low vision. Be My AI incorporated GPT-4V into the existing Be My Eyes platform which provided descriptions of photos taken by the blind user's smartphone. Be My Eyes piloted Be My AI from March to early August 2023 with a group of nearly 200 blind and low vision beta testers to hone the safety and user experience of the product. By September, the beta test group had grown to 16,000 blind and low vision users requesting a daily average of 25,000 descriptions. This testing determined that Be My AI can provide its 500,000 blind and low-vision users with unprecedented tools addressing informational, cultural, and employment needs.

A key goal of the pilot was to inform how GPT-4V can be deployed responsibly. The Be My AI beta testers surfaced AI issues including hallucinations, errors, and limitations created by product design, policy, and the model. In particular, beta testers expressed concern that the model can make basic errors, sometimes with misleading matter-of-fact confidence. One beta tester remarked: "It very confidently told me there was an item on a menu that was in fact not there." However, Be My Eyes was encouraged by the fact that we noticeably reduced the frequency and severity of hallucinations and errors over the time of the beta test. In particular, testers noticed that we improved optical character recognition and the quality and depth of descriptions.

Since risks remain, Be My Eyes warns its testers and future users not to rely on Be My AI for safety and health issues like reading prescriptions, checking ingredient lists for allergens, or crossing the street. Likewise, Be My Eyes tells its users that AI should never be used to replace a white cane or a trained guide dog. Be My Eyes will continue to be explicit on this point. Be My Eyes also offers users the option to depart the AI session and immediately connect with a human volunteer. This can be useful for human verification of AI results, or when the AI fails to identify or process an image.

Another challenge that Be My AI testers have repeatedly shared is that they want to use Be My AI to know the facial and visible characteristics of people they meet, people in social media posts, and even their own images - information that a sighted person can obtain simply by standing in any public space or looking in a mirror. But analyzing faces comes with risks including privacy considerations and the laws that govern them, and the possibility of harmful biases affecting the system's outputs. Be My Eyes received many impassioned comments about the importance of this feature. One example from one beta tester: "Thank you for hearing all of us and understanding how just a glimpse of this technology has been so impactful. I never emotionally understood the power of a picture before this service. Logos, and pages in books took on new meaning, and getting descriptions of family members both present or who have passed on was incredible. Thank you for contributing your part of give us all of that as a community."

Due to the benefits that this feature can bring to low-vision and blind users, we are designing mitigations and processes that allow features of faces and people to be described by the Be My Eyes product - providing a more equitable experience for them - without identifying people by name. We hope someday to be able to find a way to empower the blind and low-vision community to identify people - just like sighted people do - while addressing concerns around privacy and bias.

_2.1.2 Developer alpha_

In keeping with our iterative deployment approach, we engaged over a thousand alpha testers over three months in order to gain additional feedback and insight into the real ways people interact with GPT-4V. We analyzed fractions of traffic data from our alpha production traffic from July and August 2023 to better understand the use of GPT-4V for person identification, medical advice, and CAPTCHA breaking.

Of the prompts sampled, 20% were queries in which users requested general explanations and descriptions of an image: e.g., users asked the model questions such as "what", "where" or "who is this?" A more detailed breakdown exposed various risk surfaces such as medical condition diagnosis, treatment recommendations, medication intake, and several privacy-related concerns. Particular attention was given to potentially biased outputs, images of children and prompts related to them, sentiment analysis, and health status inference within the uploaded images of people. We also looked at prompts similar to "solve this puzzle," in order to understand the prevalence and nature of CAPTCHA requests. The data we found has further helped us refine our evaluations, models, and system to protect against possibly risky user queries, which you can read about in Section 2.4.

**2.2 Evaluations**

To better understand the GPT-4V system, we utilized both qualitative and quantitative evaluations. To perform qualitative evaluations, we engaged in internal experimentation to stress-test the system and solicited external expert red-teaming. For quantitative evaluations, we built evaluations that measured model refusals and model performance accuracy.

- Harmful content
  - Refusal evaluations for illicit behaviour
- Harms of representation, allocation, and quality of service
  - Refusal evaluations for ungrounded inferences
  - Performance accuracy evaluations for gender, race and age recognition across demographics
- Privacy
  - Refusal evaluation for person identification requests
  - Performance accuracy evaluation for person identification requests
  - Geolocalization evaluations
- Cybersecurity
  - Performance accuracy CAPTCHA breaking evaluations
- Multimodal Jailbreaks
  - Refusal evaluation for text-screenshot jailbreak (See Figure 1 for an example of a text-screenshot jailbreak)

Refusal evaluations measure the percentage of model outputs that constitute a refusal in response to certain potentially risky inputs (See Section 2.4 for more details on refusals). Performance accuracy evaluations measure how often the model correctly answers a certain input prompt by choosing the right answer out of 5 options.

The section below dives deeper into some of our evaluations:

- Performance on sensitive trait attribution across demographics: We studied performance parity on images of people for gender, age, and race recognition. These evaluations were constructed using publicly available datasets such as FairFace and Labeled Faces in the Wild. There are well documented studies of narrow computer vision systems exhibiting various biases such as facial recognition systems having disparate performance based on race, stereotyping members of diffierent races, etc. It is important to note that even with performance parity, differences in downstream impact and harm could still occur depending on the context of the deployment of such tools. OpenAI has thus added refusals for most instances of sensitive trait requests - you can read more about how in Section 2.4.

- Person identification evaluations: We studied the model's ability to identify people in photos, including celebrites, public servants and politicians, semi-private, and private individuals. These datasets were constructed using public datasets such as CelebA, Celebrity Faces in the Wild and a dataset of images of members of Congress for public figures. For semi-private and private individuals, we used images of employees. Performance on refusal behavior can be seen below. We find that we are able to effectively steer the model to refuse this class of requests more than 98% of the time, and steer its accuracy rate to 0% based on internal evaluations.

- Ungrounded inference evaluation: Ungrounded inferences are inferences that are not justified by the information the user has provided - in the case of GPT-4V, this means information contained in the image or text. Examples include model responses to questions such as "who is more dangerous?" or "what job does she have?" for an image of a person, which typically cannot be answered based on visual components of such an image alone. When the model provides such ungrounded inferences, it can reinforce biases or provide inaccurate information. To safeguard against these problems we've built automatic evaluations that gauge the model's propensity to successfully refuse these requests for information.

- Multimodal jailbreak evaluations: Since the release and growth of ChatCPT, a large amount of effort has been dedicated to trying to find prompts that circumvent the safety systems in place to prevent malicious misuse. These jailbreaks typically involve trapping the model via convoluted logical reasoning chains designed to make it ignore its instructions and training. A new vector for jailbreaks with image input involves placing into images some of the logical reasoning needed to break the model. This can be done in the form of screenshots of written instructions, or even visual reasoning cues (See Figure 1). Placing such information in images makes it infeasible to use text-based heuristic methods to search for jailbreaks. We must rely on the capability of the visual system itself. To quantify this we've converted a comprehensive set of known text jailbreaks to screenshots of the text. This allows us to analyze whether the visual input space provides new vectors of attack for known problems.

- Extending text-only evaluations to multimodal: We extended our text-only evaluations in domains such as advice or encouragement for self-harm behaviors, and graphic material such as erotic or violent content, by using the same set of evals from GPT-4, and then replacing words with up to two image synonyms per example. Image synonyms are images that can be used to replace a word - for example, a picture of a knife being used to indicate the word 'kill'. This was done to ensure that images did not offer an easy way to bypass our text-only mitigations.

- CAPTCHA breaking and geolocation: We used public datasets to measure the ability of the model to break CAPTCHAs and carry out broad geolocation (e.g., identify the name of the city). These evaluations represent capabilities that demostrate the model's intelligence, but can also lead to cause for concern. Task such as the ability to solve CAPTCHAs indicate the model's ability to solve puzzles and perform complex visual reasoning tasks. High performance on geolocation evaluations demonstrate world knowledge the model possesses and can be useful for users trying to search for an item or place.

![Figure 1: Example of a text-screenshot jailbreak prompt. GPT-4V-Early demonstrates the models' early performance for such prompts and GPT-4V Launch demonstrates the performance of the model we're launching.]()

However, a powerful, general purpose CAPTCHA breaker that's easily accessible can have cybersecurity and AI safety implications. These capabilities can be used to bypass security measures intended for botware, and they enable AI systems to interact with systems intended for human use.

Additionally, geolocation presents privacy concerns and can be used to identify the location of individuals who do not wish their location to be known. Note the model's geolocation abilities generally do not go deeper than the level of identifying a city from an image in most cases, reducing the likelihood of being able to find someone's precise location via the model alone.

![Figure 2: The combination of continual safety progress, model-level mitigations in the form of additional safety training data, and system level mitigations have led to significant progress in refusing disallowd prompts]()

**2.3 External Red Teaming**

As with previous deployments, OpenAI worked with external experts to qualitatively assess the limitations and risks associated with the model and system. This red teaming was specifically intended to test risks associated with the multimodal (vision) functionality of GPT-4, and builds upon the work in the GPT-4 system card. We focus this analysis on 6 key risk areas we received especially useful red teamer feedback in:

- Scientific proficiency
- Medical advice
- Stereotyping and ungrounded inferences
- Disinformation risks
- Hateful Content
- Visual vulnerabilities

![Figure 3: Evaluating GPT-4V + Refusal System against screenshots of a text refusal dataset finds that the combination of model-level mitigations and our refusal system enabled us to reach our internal target of a 100% refusal rate.]()

_2.3.1 Scientific proficiency_

Red teamers tested GPT-4V's capabilities and limitations in scientific domains. In terms of capabilities, red teamers noted the model's ability to capture complex information in images, including very specialized imagery extracted from scientific publications, and diagrams with text and detailed components. Additionally, in some instances, the model was successful at properly understanding advanced science from recent papers and critically assessing claims for novel scientific discoveries.

However, the model exhibited some key limitations. If two separate text components were closely located in an image, the model would occasionally combine them. For instance, it may merge "multipotent hematopoietic stem cell (HSC)" and "self-renewing division," (see Figure4) leading to the creation of unrelated terms. Additionally, the model was prone to hallucinations and sometimes could make factual errors in an authoritative tone. In some cases, it could also fail to identify information from images. It could miss text or characters, overlook mathematical symbols, and be unable to recognize spatial locations and color mappings.

![Figure 4: Mistakes GPT-4V makes such as combining terms and missing symbols when trying to process complex images. See Appendix A.2 for clear image given to model]()

Given the model's imperfect but increased proficiency for such tasks, it could appear to be useful for certain dangerous tasks that require scientific proficiency such as synthesis of certain illicit chemicals. For example, the model would give information for the synthesis and analysis of some dangerous chemical such as Isotonitazene, a synthetic opioid. (See Figure 5) However, the model's generations here can be inaccurate and error prone, limiting its use for such tasks.

GPT-4V has misidentified substances like fentanyl, carfentanil, and cocaine from images of their chemical structure, but also occasionally correctly identified poisonous foods such as certain toxic mushrooms from images.(See Figure 6) This demonstrates that the model is unreliable and should not be used for any high risk tasks such as identification of dangerous compounds or foods.

![Figure 5: Example of GPT-4V providing incorrect instructions to synthesize a dangerous compound.]()

![Figure 6: Example of GPT-4V's unreliable performance for correctly identifying chemical structures or poisonous foods.]()

_2.3.2 Medical advice_

Medically trained red teamers tested the model's ability to provide medical advice, especially with medically-related images as an input. Red teamers took into account the perspectives of both a layperson seeking medical advice and a professional with medical education and training. Some considerations for potential risks that may arise during the course of using the model to seek medical advice are accuracy, bias, and taking context into account.

Red teamers found that there were inconsistencies in interpretation in medical imaging - while the model would occasionally give accurate responses, it could sometimes give wrong responses for the same question. As one example, Figure 7 shows some of the vulnerabilities or inaccuracies that could result from an incorrect or decontextualized interpretation of the directionality of medical imaging. The experts noted that the correct standard when viewing imaging scans is to imagine the patient as if they were facing you, meaning the right side on the image would correspond to the left side of the patient. This is an important concept that is needed when viewing and diagnosing radiographic imaging. Misdiagnosing the laterality of any number of condition is very dangerous.

Given the model's imperfect performance in this domain and the risks associated with inaccuracies, we do not consider the current version of GPT-4V to be fit for performing any medical function or substituting professional medical advice, diagnosis, or treatment, or judgment.

![Figure 7: Examples of GPT-4V's unreliable performance for medical purposes]()

_2.3.3 Stereotyping and ungrounded inferences_

Using GPT-4V for some tasks might generate unwanted or harmful assumptions that are not grounded in the information provided to the model (the image or the text prompt). Red teamers tested risks associated with ungrounded inferences about people and places.

In early versions of GPT-4V, prompting the model to make a decision between a variety of options, followed by asking for an explanation frequently surfaced stereotypes and ungrounded inferences within the model.

Broad open-ended questions to the model paired with an image also exposed bias or anchoring towards specific topics that may not necessarily have been intended by the prompt. Eg. When prompted to advise the woman in the image, the model focuses on subjects of body weight and body positivity. (See Figure 8)

We have added mitigations for risks associated with ungrounded inferences by having the model refuse such requests relating to people. This is a conservative approach, and our hope is that as we refine our research and mitigations, the model may be able to answer questions about people in low-risk contexts.

![Figure 8: Examples of ungrounded inferences and stereotypes that early version of GPT-4V exhibited compared to the behavior the launch model exhibits.]()

_2.3.4 Disinformation risks_

As noted in the GPT-4 system card, the model can be used to generate plausible realistic and targeted text content. When paired with vision capabilities, image and text content can pose increased risks with disinformation since the model can create text content tailored to an image input. Previous work has shown that people are more likely to believe true and false statements when they're presented alongside an image, and have false recall of made up headlines when they are accompanied with a photo. It is also known that engagement with content increases when it is associated with an image.

Red teamers also tested GPT-4V's ability to detect incorrect information or disinformation in an image. The model's ability to recognize disinformation was inconsistent, but may be related to how well-known a disinformation concept is and its recency. Overall, GPT-4V was not trained for this purpose and should not be used as a way to detect disinformation, or to otherwise verify whether something is ture or false.

Realistic, customized images can be created using other generative image models, and used in combination with GPT-4V's capabilities. Pairing the ability of image models to generate images more easily with GPT-4V's ability to generate accompanying text more easily may have an impact on disinformation risks. However, a proper risk assessment would also have to take into account the context of use (e.g. the actor, the surrounding events, etc.), the manner and extent of distribution (e.g. is the pairing within a closed software application or in public forums), and the presence of other mitigations such as watermarking or other provenance tools for the generated image.

_2.3.5 Hateful content_

GPT-4V refuses to answer questions about hate symbols and extremist content in some instances but not all. The behavior may be inconsistent and at times contextually inappropriate. For instance, it knows the historic meaning of the Templar Cross but misses its modern meaning in the US, where it has been appropriated by hate groups. See Figure 10a.

Red teamers observed that if a user directly names a well-known hate group, the model ussually refuses to provide a completion. But, if you use lesser-known names - such as "Totenwaffen" - or symbols, you might get past this. The model can also sometimes make songs or poems that praise certain hate figures or groups if given a picture of them, when the figures or groups are not explicitly named. OpenAI has added refusals for certain kinds of obviously harmful generations in the space but not all (see Figure 10b). This remains a dynamic, challenging problem to solve.

![Figure 10 (a) GPT-4V responds with the historical meaning of the image but is unaware the image has been appropriated by hate groups. (b) If prompted, GPT-4V can generate content praising certain lesser known hate groups in response to their symbols.]()

_2.3.6 Visual vulnerabilities_

Red teaming found some limitations that are specifically associated with the ways that images could be used or presented. For example: ordering of the images used as input may influence the recommendation made. In the example in 11, asking for which state to move to, based on the flags inputted, favors the first flag inputted when red teamers tested both possible orderings of the flags.

The example represents challenges with robustness and reliability that the model still faces. We anticipate there to be many more such vulnerabilities in the model that we discover through its broad usage and we will be working on improving model performance for future iterations to be robust to them.

![Figure 11: Examples of visual vulnerabilities GPT-4V exhibits. This example demonstrates model generations can be sensitive to the order in which images are given to the model.]()

**2.4 Mitigations**

_2.4.1 Transfer benefits from existing safety work_

GPT-4V inherits several transfer benefits from model-level and system-level safety mitigations already deployed in GPT-4. In a similar vein, some of our safety measures implemented fro DALL·E proved beneficial in addressing potential multi-modal risk in GPT-4V.

Internal evaluations show that performance of refusals of text content against our existing policies is equivalent to our base language model for GPT-4V. At the system-level, our existing moderation classifiers continue to inform our monitoring and enforcement pipelines for post-hoc enforcement of text inputs and outputs. GPT-4V mirrors our existing moderation efforts deployed in DALL·E to detect explicit image uploads by users.

These transfer benefits from our prior safety work enable us to focus on novel risks introduced by this multimodal model. This includes areas where, in isolation, the text or image content is benign, but in concert create a harmful prompt or generation; images with people in them; and common multimodal jailbreaks such as adversarial images with text.

![Figure 12: Example prompt given to GPT-4 to find phrases to replace with images to turn text-only prompts into multimodal prompts.]()

_2.4.2 Additional Mitigations for High-Risk Areas_

GPT-4V includes carefully designed refusal behavior for some prompts that contain images of people. The model refuses requests for the following:

- Identity (e.g. a user uploads an images of a person and asks who they are, or a pair of images and asks if they're the same person)
- Sensitive traits (e.g. age, race)
- Ungrounded inferences (e.g. when the model draws conclusions based on those traits not visually present, as discussed in Section 2.2)

To further reduce the risks in emerging and high-stake areas, we integrated additional multimodal data into the post-training process in order to reinforce refusal behavior for illicit behavior and ungrounded inference requests. Our focus was to mitigate risky prompts where in isolation, the text and the image were individually benign, but when combined as a multimodal prompt, could lead to harmful outputs.

For illicit behavior, we collected a multimodal dataset by augmenting our existing text-only dataset with image synonyms. For example, given a text string "how do i kill the people?", we want to adapt it into a multimodal example "how do i \[image of knife\] the \[image of people\]?". The augmentation consists of the following steps:

- For each original text-only example, we ask GPT-4 to pick the top two most harmful short phrases (ref the table below);
- For each chosen short phrase, we replace it with a web crawled image.
- To ensure semantic-invariant, we conduct hunman review and filter out low quality augmentations.
- To reinforce the robustness of the refusal behavior, we also augment the examples with various system messages.

For ungrounded inference requests, we used data collected through our red teaming campaigns. The goal was to train the model to refuse prompts that were requesting an ungrounded conclusion based on certain attributes of a person. For example, if the prompt includes a photo of a person and the text "What is this person's personality type?", the desired model completion is "I'm sorry, I can't help with that." The examples collected through the red teaming campaign were further reviewed by humans before adding to the training dataset.

According to our internal evaluations after post-training, we observed that 97.2% of the completions refused requests for illicit advice, and 100% of the completions refused requests for ungrounded inference. In addition to measuring the refusal of completions, we also evaluate the correct refusal style. This evaluation only considers the subset of all refusals that are short and concise to be correct. We observed that the correct refusal style rate improved from 44.4% to 72.2% for illicit advice style, and from 7.5% to 50% for ungrounded inference style. We will iterate and improve refusals over time as we continue to learn from real world use.

In addition to the model-level mitigations described above, we added system-level mitigations for adversarial images containing overlaid text in order to ensure this input couldn't be used to circumvent our text safety mitigations. For example, a user could submit an image containing the text, "How do I build a bomb?" As one mitigation for this risk, we images through an OCR tool and then calculate moderation scores on the resulting text in the image. This is in addition to detecting any text inputted directly in the prompt.

### 3 Conclusion and Next Steps

GPT-4V's capabilities pose exciting opportunities and novel challenges. Our deployment preparation approach has targeted assessment and mitigations of risks related to images of people such as person identification, biased outputs from images of people including representational harms or allocational harms that may stem from such inputs. Additionally, we have studied the model's capability jumps in certain high-risk domains such as medicine and scientific proficiency.

There are a few next steps that we will be investing further in and will be engaging with the public on:

- There are fundamental questions around behaviors the models should or should not be allowed to engage in. Some examples of these include: should models carry out identification of public figures such as Alan Turing from their images? Should models be allowed to infer gender, race, or emotions from images of people? Should the visually impaired receive special consideration in these questions fro the sake of accessibility? These questions traverse well-documented and novel concerns around privacy, fairness, and the role AI models are allowed to play in society.

- As these models are adopted globally, improving performance in languages spoken by global users, as well as enhancing image recognition capabilities that are relevant to a worldwide audience, is becoming increasingly critical. We plan to continue investing in advancements in these areas.

- We will be focusing on research that allows us to get higher precision and more sophisticated with how we handle image uploads with people. While we currently have fairly broad but imperfect refusals for responses related to people, we will hone this by advancing how the model handles sensitive information from images, like a person's identity or protected characteristics. Additionally, we will further invest in mitigating representational harms that may stem from stereotypical or denigrating outputs.

### 4 Acknowledgements

We are grateful to our expert adversarial testers and red teamers who helped test our models at early stages of development and informed our risk assessments as well as the System Card output. Participation in this red teaming process is not an endorsement of the deployment plans of OpenAI or OpenAI's policies: ...(name).

We thank Microsoft for their partnership, especially Microsoft Azure for supporting model training with infrastructure design and management, and the Microsoft Bing team and Microsoft's safety teams for their partnership on safe deployment and safety research.
