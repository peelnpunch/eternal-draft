# "The Eternal Draft" Message Dispatcher 

🔥🌌 Welcome to the "The Eternal Draft" message dispatcher - a Go application crafted with love by the inhabitants of "The Eternal Draft" at AfrikaBurn. ♾️ 📐🌌🔥

We are a collective of makers, builders, and good time havers. 🛠️🕺🏼  The name 'The Eternal Draft' embodies the philosophy that all creation and our dasein is an ongoing process of iteration.

Our centerpiece? A portal to message your future self, embracing the notion that we are all works in progress. 📩➡️👤 Here, creation isn't just a vibe—it's a continuous state of doing; a product of diligence and hard work, no matter the direction you take it.

🎨💪 We'll offer folks on their way to/from the playa to join us in our temporary enclave and scream into our portal, before we all return to the drafting table. 🏜️🗣️

## Features

- **SMTP Configuration**: Easily configure SMTP server settings to ensure seamless email delivery.
- **Emails with Attachments**: Supports sending both textual content and image attachments in emails.
- **Automated Payload Parsing**: Automatically parses email payload information from filenames within a specified directory.
- **Time-based Email Triggering**: Sends emails based on the calculated difference between the current date and a pre-set date, facilitating scheduled email dispatches.

## Prerequisites

Before you begin, ensure you have met the following requirements:
- Go version 1.15 or newer installed on your machine.
- Valid credentials and access to an SMTP server for email dispatch.

## Installation

To install the Eternal Draft Email Sender, follow these steps:

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/peelnpunch/eternal-draft
    ```

2. Navigate to the project directory:

    ```bash
    cd eternal-draft
    ```

## Configuration

1. **SMTP Settings**: Locate the `SMTPSettings` struct initialization in the `main` function and update it with your SMTP server details.

2. **Postcards Directory**: Ensure that your postcard images are placed within a directory named `postcards` at the root of your project. The file naming should adhere to the `email_years.extension` format.

## Running the Application

Execute the application with the following command:

```bash
go run .
