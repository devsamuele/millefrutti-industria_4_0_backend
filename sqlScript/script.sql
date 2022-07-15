USE [ADB_MILLEFRUTTISRL]
GO

SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[xPastorizzatore]
(
	[id] [int] IDENTITY(1,1) NOT NULL,
	[cd_lotto] [varchar](20) NOT NULL,
	[cd_ar] [varchar](20) NOT NULL,
	[basil_amount] [int] NOT NULL,
	[elaborazione] [bit] DEFAULT 0 NULL,
	[packages] [int] NOT NULL,
	[date] [datetime] NOT NULL,
	[document_created] [bit] NOT NULL,
	[status] [varchar](255) NOT NULL,
	[created] [datetime] NOT NULL,
	CONSTRAINT [PK_xPastorizzatore] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY],
	CONSTRAINT [cd_ar_xPastorizzatore] UNIQUE NONCLUSTERED 
(
	[cd_ar] ASC,
	[cd_lotto] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO

SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[xCentrifuga]
(
	[id] [int] IDENTITY(1,1) NOT NULL,
	[cd_lotto] [varchar](20) NOT NULL,
	[cd_ar] [varchar](20) NOT NULL,
	[cycles] [int] NOT NULL,
	[elaborazione] [bit] DEFAULT 0 NULL,
	[total_cycles] [int] NOT NULL,
	[date] [datetime] NOT NULL,
	[document_created] [bit] NOT NULL,
	[status] [varchar](255) NOT NULL,
	[created] [datetime] NOT NULL,
	CONSTRAINT [PK_xCentrifuga] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY],
	CONSTRAINT [cd_ar_xCentrifuga] UNIQUE NONCLUSTERED 
(
	[cd_ar] ASC,
	[cd_lotto] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO

EXEC asp_du_AddAlterColumn 'Dotes', 'xId_Centrifuga', 'int NULL', '', 'ID di xCentrifuga'
EXEC asp_du_AddAlterColumn 'Dotes', 'xId_Pastorizzatore', 'int NULL', '', 'ID di xPastorizzatore'

