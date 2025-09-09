export default ({ env }) => {
	const awsBucket = env("AWS_BUCKET");
	const awsRegion = env("AWS_REGION", "ap-southeast-1");
	const s3AssetsUrl = awsBucket
		? `${awsBucket}.s3.${awsRegion}.amazonaws.com`
		: "";

	return {
		email: {
			config: {
				provider: "nodemailer",
				providerOptions: {
					host: env("SMTP_HOST", "localhost"),
					port: env("SMTP_PORT", 1025),
					auth: {
						user: env("SMTP_USERNAME"),
						pass: env("SMTP_PASSWORD"),
					},
					ignoreTLS: process.env.NODE_ENV === "development",
				},
				settings: {
					defaultFrom: env("SMTP_EMAIL_FROM", "admin@example.com"),
					defaultReplyTo: env("SMTP_REPLY_TO", "admin@example.com"),
				},
			},
		},
		upload: {
			config: {
				provider: "aws-s3",
				providerOptions: {
					baseUrl: s3AssetsUrl,
					rootPath: env("S3_PATH_PREFIX"),
					s3Options: {
						credentials: {
							accessKeyId: env("AWS_ACCESS_KEY_ID"),
							secretAccessKey: env("AWS_ACCESS_SECRET"),
						},
						region: awsRegion,
						params: {
							ACL: env("AWS_ACL", "public-read"),
							signedUrlExpires: env("AWS_SIGNED_URL_EXPIRES", 15 * 60),
							Bucket: env("AWS_BUCKET"),
						},
					},
				},
				actionOptions: {
					upload: {},
					uploadStream: {},
					delete: {},
				},
			},
		},
	};
};
